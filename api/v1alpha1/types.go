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

package v1alpha1

// ResourceLifecycle configures the lifecycle of a resource.
type ResourceLifecycle string

// Region is a Hivelocity Location
// +kubebuilder:validation:Enum=AMS1;ATL2;BOM1;DAL1;DEL1;EDGE-ARN1;EDGE-CDG1;EDGE-FLL1;EDGE-FRA1;EDGE-HKG1;EDGE-IAD1;EDGE-ICN1;EDGE-JFK1;EDGE-LAX1;EDGE-LCY1;EDGE-LIN1;EDGE-NRT1;EDGE-SIN1;EDGE-SNV1;EDGE-SYD1;EDGE-TOJ1;EDGE-YXX1;EDGE-YYZ1;FRA1;IAD3;IND1;LAX2;LHR2;MIA1;NRT2;NYC1;ORD1;PNQ1;POZ1;RIX1;SEA1;SIN1;SLC1;TPA1;TPA2;VNO1;YYZ2
type Region string
