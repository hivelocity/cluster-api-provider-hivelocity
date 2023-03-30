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
	"time"

	infrav1 "github.com/hivelocity/cluster-api-provider-hivelocity/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// actionResult is an interface that encapsulates the result of a Reconcile
// call, as returned by the action corresponding to the current state.
type actionResult interface {
	Result() (reconcile.Result, error)
}

// actionContinue is a result indicating that the current action is still
// in progress, and that the resource should remain in the same provisioning
// state.
type actionContinue struct {
	delay time.Duration
}

func (r actionContinue) Result() (result reconcile.Result, err error) {
	result.RequeueAfter = r.delay
	// Set Requeue true as well as RequeueAfter in case the delay is 0.
	result.Requeue = true
	return
}

// actionComplete is a result indicating that the current action has completed,
// and that the resource should transition to the next state.
type actionComplete struct{}

func (r actionComplete) Result() (result reconcile.Result, err error) {
	result.Requeue = true
	return
}

// actionError is a result indicating that an error occurred while attempting
// to advance the current action, and that reconciliation should be retried.
type actionError struct {
	err error
}

func (r actionError) Result() (result reconcile.Result, err error) {
	err = r.err
	return
}

// actionGoBack is a result indicating that the current action cannot be carried out
// and that the resource should transition to a previous state.
type actionGoBack struct {
	nextState infrav1.ProvisioningState
}

func (r actionGoBack) Result() (result reconcile.Result, err error) {
	result.Requeue = true
	return
}

// actionFailed is a result indicating that the current action has failed,
// and that the resource should be marked as in error.
type actionFailed struct{}

func (r actionFailed) Result() (result reconcile.Result, err error) {
	result.RequeueAfter = 1 * time.Minute
	return
}
