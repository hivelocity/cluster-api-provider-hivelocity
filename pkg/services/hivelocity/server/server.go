/*
Copyright 2023 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package server implements functions to manage the lifecycle of Hivelocity servers.
package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	infrav1 "github.com/hivelocity/cluster-api-provider-hivelocity/api/v1alpha1"
	"github.com/hivelocity/cluster-api-provider-hivelocity/pkg/scope"
	hvclient "github.com/hivelocity/cluster-api-provider-hivelocity/pkg/services/hivelocity/client"
	hv "github.com/hivelocity/hivelocity-client-go/client"
	corev1 "k8s.io/api/core/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	capierrors "sigs.k8s.io/cluster-api/errors"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const maxShutDownTime = 2 * time.Minute
const serverOffTimeout = 10 * time.Minute
const ErrorCodeRateLimitExceeded = 429

// Service defines struct with machine scope to reconcile Hivelocity machines.
type Service struct {
	scope *scope.MachineScope
}

// NewService outs a new service with machine scope.
func NewService(scope *scope.MachineScope) *Service {
	return &Service{
		scope: scope,
	}
}

// Reconcile implements reconcilement of Hivelocity machines.
func (s *Service) Reconcile(ctx context.Context) (_ *ctrl.Result, err error) {
	if s.scope.HivelocityCluster.Spec.HivelocitySecret.Key == "" {
		record.Eventf(s.scope.HivelocityMachine, corev1.EventTypeWarning, "NoAPIKey", "No Hivelocity API Key found")
		return nil, fmt.Errorf("no Hivelocity API Key provided - cannot reconcile Hivelocity server")
	}

	// detect failure domain
	failureDomain, err := s.scope.GetFailureDomain()
	if err != nil {
		return nil, fmt.Errorf("failed to get failure domain: %w", err)
	}
	s.scope.HivelocityMachine.Status.Region = infrav1.Region(failureDomain)

	// Waiting for bootstrap data to be ready
	if !s.scope.IsBootstrapDataReady(ctx) {
		s.scope.Info("Bootstrap not ready - requeuing")
		conditions.MarkFalse(
			s.scope.HivelocityMachine,
			infrav1.InstanceBootstrapReadyCondition,
			infrav1.InstanceBootstrapNotReadyReason,
			clusterv1.ConditionSeverityInfo,
			"bootstrap not ready yet",
		)
		return &ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}

	conditions.MarkTrue(
		s.scope.HivelocityMachine,
		infrav1.InstanceBootstrapReadyCondition,
	)

	// Try to find an existing server
	server, err := s.findServer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get server: %w", err)
	}

	// If no server is found we have to create one
	if server == nil {
		server, err = s.createServer(ctx, failureDomain)
		if err != nil {
			return nil, fmt.Errorf("failed to create server: %w", err)
		}
		record.Eventf(
			s.scope.HivelocityMachine,
			"SuccessfulCreate",
			"Created new server with id %d",
			server.DeviceId,
		)
	}

	c := s.scope.HivelocityMachine.Status.Conditions.DeepCopy()
	s.scope.HivelocityMachine.Status = setStatusFromAPI(server)
	s.scope.HivelocityMachine.Status.Conditions = c

	switch server.PowerStatus {
	case hvclient.PowerStatusOff:
		return s.handleServerStatusOff(ctx, server)
	/* todo: up to now we only get ON/OFF
	case hv.BareMetalDeviceStatusStarting:
		// Requeue here so that server does not switch back and forth between off and starting.
		// If we don't return here, the condition InstanceReady would get marked as true in this
		// case. However, if the server is stuck and does not power on, we should not mark the
		// condition InstanceReady as true to be able to remediate the server after a timeout.
		return &reconcile.Result{RequeueAfter: 10 * time.Second}, nil
	*/
	case hvclient.PowerStatusOn: // Do nothing
	default:
		s.scope.HivelocityMachine.Status.Ready = false
		s.scope.V(1).Info("server not in running state",
			"server", server.Hostname,
			"powerStatus", server.PowerStatus)
		return &reconcile.Result{RequeueAfter: 2 * time.Second}, nil
	}

	providerID := fmt.Sprintf("hivelocity://%d", server.DeviceId)

	if !s.scope.IsControlPlane() {
		s.scope.HivelocityMachine.Spec.ProviderID = &providerID
		s.scope.HivelocityMachine.Status.Ready = true
		conditions.MarkTrue(s.scope.HivelocityMachine, infrav1.InstanceReadyCondition)
		return nil, nil
	}

	/*
		// all control planes have to be attached to the load balancer if it exists
		if err := s.reconcileLoadBalancerAttachment(ctx, server); err != nil {
			return nil, fmt.Errorf("failed to reconcile load balancer attachement: %w", err)
		}
	*/

	s.scope.HivelocityMachine.Spec.ProviderID = &providerID
	s.scope.HivelocityMachine.Status.Ready = true
	conditions.MarkTrue(s.scope.HivelocityMachine, infrav1.InstanceReadyCondition)

	return nil, nil
}

func (s *Service) handleServerStatusOff(ctx context.Context, server *hv.BareMetalDevice) (*reconcile.Result, error) {
	// Check if server is in ServerStatusOff and turn it on. This is to avoid a bug of Hivelocity where
	// sometimes machines are created and not turned on

	condition := conditions.Get(s.scope.HivelocityMachine, infrav1.InstanceReadyCondition)
	if condition != nil &&
		condition.Status == corev1.ConditionFalse &&
		condition.Reason == infrav1.ServerOffReason {
		if time.Now().Before(condition.LastTransitionTime.Time.Add(serverOffTimeout)) {
			// Not yet timed out, try again to power on
			if err := s.scope.HVClient.PowerOnServer(ctx, server); err != nil {
				return nil, fmt.Errorf("failed to power on server: %w", err)
			}
		} else {
			// Timed out. Set failure reason
			s.scope.SetError("reached timeout of waiting for machines that are switched off", capierrors.CreateMachineError)
			return nil, nil
		}
	} else { // todo: too complicated. Is there a way to avoid the "else"?
		// No condition set yet. Try to power server on.
		if err := s.scope.HVClient.PowerOnServer(ctx, server); err != nil {
			if errors.Is(err, hvclient.ErrRateLimitExceeded) {
				conditions.MarkTrue(s.scope.HivelocityMachine, infrav1.RateLimitExceeded)
				record.Event(s.scope.HivelocityMachine,
					"RateLimitExceeded",
					"exceeded rate limit with calling hivelocity function PowerOnServer",
				)
			}
			return nil, fmt.Errorf("failed to power on server: %w", err)
		}
		conditions.MarkFalse(
			s.scope.HivelocityMachine,
			infrav1.InstanceReadyCondition,
			infrav1.ServerOffReason,
			clusterv1.ConditionSeverityInfo,
			"server is switched off",
		)
	}

	// Try again in 30 sec.
	return &reconcile.Result{RequeueAfter: 30 * time.Second}, nil
}

func (s *Service) reconcileNetworkAttachment(ctx context.Context, server *hv.BareMetalDevice) error {
	return nil
}

func (s *Service) createServer(ctx context.Context, failureDomain string) (*hv.BareMetalDevice, error) {
	// get userData
	userData, err := s.scope.GetRawBootstrapData(ctx)
	if err != nil {
		record.Warnf(
			s.scope.HivelocityMachine,
			"FailedGetBootstrapData",
			err.Error(),
		)
		return nil, fmt.Errorf("failed to get raw bootstrap data: %s", err)
	}

	image, err := s.getServerImage(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get server image: %w", err)
	}

	opts := hv.BareMetalDeviceUpdate{
		Hostname: s.scope.Name(),
		//todo Tags: createLabels(s.scope.HivelocityCluster.Name, s.scope.Name(), s.scope.IsControlPlane()),
		Script: string(userData), // cloud-init script
		OsName: image,
	}

	if s.scope.HivelocityCluster.Spec.SSHKey != nil {
		sshKeyID, err := getSSHKeyIDFromSSHKeyName(s.scope.HivelocityCluster.Spec.SSHKey)
		if err != nil {
			return nil, fmt.Errorf("error with ssh keys: %w", err)
		}

		opts.PublicSshKeyId = sshKeyID
	}

	// Create the server
	server, err := s.scope.HVClient.CreateServer(ctx, &opts)
	if err != nil {
		if errors.Is(err, hvclient.ErrRateLimitExceeded) {
			conditions.MarkTrue(s.scope.HivelocityMachine, infrav1.RateLimitExceeded)
			record.Event(s.scope.HivelocityMachine,
				"RateLimitExceeded",
				"exceeded rate limit with calling Hivelocity function CreateServer",
			)
		}
		record.Warnf(s.scope.HivelocityMachine,
			"FailedCreateHivelocityServer",
			"Failed to create Hivelocity server %s: %s",
			s.scope.Name(),
			err,
		)
		return nil, fmt.Errorf("error while creating Hivelocity server %s: %s", s.scope.HivelocityMachine.Name, err)
	}

	return server, nil
}

func (s *Service) getServerImage(ctx context.Context) (string, error) {
	return "todo", fmt.Errorf("todo: getServerImage()")
}

func getSSHKeyIDFromSSHKeyName(sshKeyID *infrav1.SSHKey) (int32, error) {
	if sshKeyID == nil {
		return 0, fmt.Errorf("SSHKey is nil")
	}
	// todo: we only support one ssh-key for all nodes of the cluster.
	return 0, fmt.Errorf("todo: getSSHKeyIDFromSSHKeyName()")
}

// Delete implements delete method of server.
func (s *Service) Delete(ctx context.Context) (_ *ctrl.Result, err error) {
	return nil, fmt.Errorf("todo: Delete()")
}

func (s *Service) handleDeleteServerStatusRunning(ctx context.Context, server *hv.BareMetalDevice) (*ctrl.Result, error) {
	return nil, fmt.Errorf("todo: handleDeleteServerStatusRunning()")
}

func (s *Service) handleDeleteServerStatusOff(ctx context.Context, server *hv.BareMetalDevice) (*ctrl.Result, error) {
	return nil, fmt.Errorf("todo: handleDeleteServerStatusOff()")
}

func setStatusFromAPI(server *hv.BareMetalDevice) infrav1.HivelocityMachineStatus {
	var status infrav1.HivelocityMachineStatus
	return status
}

func (s *Service) reconcileLoadBalancerAttachment(ctx context.Context, server *hv.BareMetalDevice) error {
	return fmt.Errorf("todo: reconcileLoadBalancerAttachment()")
}

func (s *Service) deleteServerOfLoadBalancer(ctx context.Context, server *hv.BareMetalDevice) error {
	return fmt.Errorf("todo: deleteServerOfLoadBalancer()")
}

// We write the server name in the labels, so that all labels are or should be unique.
func (s *Service) findServer(ctx context.Context) (*hv.BareMetalDevice, error) {
	return nil, fmt.Errorf("todo: Service.findServer()")
	/*
		opts := hv.ServerListOpts{}
		opts.LabelSelector = utils.LabelsToLabelSelector(createLabels(s.scope.HivelocityCluster.Name, s.scope.Name(), s.scope.IsControlPlane()))
		servers, err := s.scope.HVClient.ListServers(ctx, opts)
		if err != nil {
			if errors.Is(err, hvclient.ErrRateLimitExceeded) {
				conditions.MarkTrue(s.scope.HivelocityMachine, infrav1.RateLimitExceeded)
				record.Event(s.scope.HivelocityMachine,
					"RateLimitExceeded",
					"exceeded rate limit with calling Hivelocity function ListServers",
				)
			}
			return nil, err
		}
		if len(servers) > 1 {
			record.Warnf(s.scope.HivelocityMachine,
				"MultipleInstances",
				"Found %v servers of name %s",
				len(servers),
				s.scope.Name())
			return nil, fmt.Errorf("found %v servers with name %s", len(servers), s.scope.Name())
		} else if len(servers) == 0 {
			return nil, nil
		}

		return servers[0], nil
	*/
}

func createLabels(hivelocityClusterName, hivelocityMachineName string, isControlPlane bool) map[string]string {
	return map[string]string{}
	/*
		m := map[string]string{
			infrav1.ClusterTagKey(hivelocityClusterName): string(infrav1.ResourceLifecycleOwned),
			infrav1.MachineNameTagKey:                hivelocityMachineName,
		}

		var machineType string
		if isControlPlane {
			machineType = "control_plane"
		} else {
			machineType = "worker"
		}
		m["machine_type"] = machineType
		return m
	*/
}
