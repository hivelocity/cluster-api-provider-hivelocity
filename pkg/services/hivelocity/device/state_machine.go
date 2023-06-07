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

package device

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	infrav1 "github.com/hivelocity/cluster-api-provider-hivelocity/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/cluster-api/util/record"
)

// stateMachine is a finite state machine that manages transitions between
// the states of a BareMetalhvMachine.
type stateMachine struct {
	hvMachine  *infrav1.HivelocityMachine
	reconciler *Service
	nextState  infrav1.ProvisioningState
	log        logr.Logger
}

func newStateMachine(hvMachine *infrav1.HivelocityMachine, reconciler *Service) *stateMachine {
	currentState := hvMachine.Spec.Status.ProvisioningState
	r := stateMachine{
		hvMachine:  hvMachine,
		reconciler: reconciler,
		nextState:  currentState, // Remain in current state by default
		log:        reconciler.scope.Logger.WithValues("caller", "state-machine"),
	}
	return &r
}

type stateHandler func(context.Context) actionResult

func (sm *stateMachine) handlers() map[infrav1.ProvisioningState]stateHandler {
	return map[infrav1.ProvisioningState]stateHandler{
		infrav1.StateAssociateDevice:         sm.handleAssociateDevice,
		infrav1.StateVerifyAssociate:         sm.handleVerifyAssociate,
		infrav1.StateProvisionDevice:         sm.handleProvisionDevice,
		infrav1.StateDeviceProvisioned:       sm.handleDeviceProvisioned,
		infrav1.StateDeleteDeviceDeProvision: sm.handleDeleteDeviceDeProvision,
		infrav1.StateDeleteDeviceDissociate:  sm.handleDeleteDeviceDissociate,
	}
}

func (sm *stateMachine) ReconcileState(ctx context.Context) (actionRes actionResult) {
	initialState := sm.hvMachine.Spec.Status
	defer func() {
		if sm.nextState != initialState.ProvisioningState {
			sm.log.Info("changing provisioning state", "old", initialState.ProvisioningState, "new", sm.nextState)
			sm.hvMachine.Spec.Status.ProvisioningState = sm.nextState
		}
		if diff := cmp.Diff(initialState, sm.hvMachine.Spec.Status); diff != "" {
			t := metav1.Now()
			sm.hvMachine.Spec.Status.LastUpdated = &t
			sm.log.V(1).Info("spec.status was updated", "diff", diff)
		}
	}()

	// we start with associating the device
	if initialState.ProvisioningState == infrav1.StateNone {
		initialState.ProvisioningState = infrav1.StateAssociateDevice
		sm.hvMachine.Spec.Status.ProvisioningState = infrav1.StateAssociateDevice
	}

	sm.log.V(1).Info("ReconcileState", "initialState.ProvisioningState", initialState.ProvisioningState)

	if stateHandler, found := sm.handlers()[initialState.ProvisioningState]; found {
		return stateHandler(ctx)
	}

	record.Warnf(sm.hvMachine, "NoHandlerFoundForState", "no handler found for state %q", initialState.ProvisioningState)
	return actionError{fmt.Errorf("no handler found for state %q", initialState.ProvisioningState)}
}

func (sm *stateMachine) handleAssociateDevice(ctx context.Context) actionResult {
	actResult := sm.reconciler.actionAssociateDevice(ctx)
	if _, ok := actResult.(actionComplete); ok {
		sm.nextState = infrav1.StateVerifyAssociate
	}
	return actResult
}

func (sm *stateMachine) handleVerifyAssociate(ctx context.Context) actionResult {
	actResult := sm.reconciler.actionVerifyAssociate(ctx)
	if _, ok := actResult.(actionComplete); ok {
		sm.nextState = infrav1.StateProvisionDevice
	}

	// check whether we need to associate the machine to another device
	actionGoBack, ok := actResult.(actionGoBack)
	if ok {
		sm.nextState = actionGoBack.nextState
	}
	return actResult
}

func (sm *stateMachine) handleProvisionDevice(ctx context.Context) actionResult {
	actResult := sm.reconciler.actionProvisionDevice(ctx)
	if _, ok := actResult.(actionComplete); ok {
		sm.nextState = infrav1.StateDeviceProvisioned
	}

	// check whether we need to go back to previous state
	actionGoBack, ok := actResult.(actionGoBack)
	if ok {
		sm.nextState = actionGoBack.nextState
	}
	return actResult
}

func (sm *stateMachine) handleDeviceProvisioned(ctx context.Context) actionResult {
	actResult := sm.reconciler.actionDeviceProvisioned(ctx)
	if _, ok := actResult.(actionComplete); ok {
		sm.nextState = infrav1.StateDeviceProvisioned
	}
	return actResult
}

func (sm *stateMachine) handleDeleteDeviceDeProvision(ctx context.Context) actionResult {
	actResult := sm.reconciler.actionDeleteDeviceDeProvision(ctx)
	if _, ok := actResult.(actionComplete); ok {
		sm.nextState = infrav1.StateDeleteDevice
	}
	return actResult
}

func (sm *stateMachine) handleDeleteDeviceDissociate(ctx context.Context) actionResult {
	actResult := sm.reconciler.actionDeleteDeviceDissociate(ctx)
	if _, ok := actResult.(actionComplete); ok {
		sm.nextState = infrav1.StateDeleteDevice
	}
	return actResult
}
