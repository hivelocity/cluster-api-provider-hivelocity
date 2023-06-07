/*
Copyright 2022 The Kubernetes Authors.

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

package e2e

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
)

const controlPlaneMachineCount int64 = 1
const workerMachineCount int64 = 1

var _ = Describe(fmt.Sprintf("[Basic] Testing Cluster %dx control-planes %dx worker and multi-tenancy",
	controlPlaneMachineCount, workerMachineCount),
	func() {
		ctx := context.TODO()

		Context("Running the CaphvClusterDeploymentSpec in hivelocity with the default flavor", func() {
			CaphvClusterDeploymentSpec(ctx, func() CaphvClusterDeploymentSpecInput {
				return CaphvClusterDeploymentSpecInput{
					E2EConfig:                e2eConfig,
					ClusterctlConfigPath:     clusterctlConfigPath,
					BootstrapClusterProxy:    bootstrapClusterProxy,
					ArtifactFolder:           artifactFolder,
					SkipCleanup:              skipCleanup,
					ControlPlaneMachineCount: controlPlaneMachineCount,
					WorkerMachineCount:       workerMachineCount,
					Flavor:                   "",
				}
			})
		})

	})
