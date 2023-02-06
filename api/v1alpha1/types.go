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

// LoadBalancerAlgorithmType defines the Algorithm type.
// +kubebuilder:validation:Enum=round_robin;least_connections
type LoadBalancerAlgorithmType string

const (

	// LoadBalancerAlgorithmTypeRoundRobin default for the Kubernetes Api Server loadbalancer.
	LoadBalancerAlgorithmTypeRoundRobin = LoadBalancerAlgorithmType("round_robin")

	// LoadBalancerAlgorithmTypeLeastConnections default for Loadbalancer.
	LoadBalancerAlgorithmTypeLeastConnections = LoadBalancerAlgorithmType("least_connections")
)

// LoadBalancerTargetType defines the target type.
// +kubebuilder:validation:Enum=server;ip
type LoadBalancerTargetType string

const (

	// LoadBalancerTargetTypeServer default for the Kubernetes Api Server loadbalancer.
	LoadBalancerTargetTypeServer = LoadBalancerTargetType("server")

	// LoadBalancerTargetTypeIP default for Loadbalancer.
	LoadBalancerTargetTypeIP = LoadBalancerTargetType("ip")
)

// SSHKey defines the SSHKey for Hivelocity.
type SSHKey struct {
	// Name of SSH key
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
	// Fingerprint of SSH key - added by controller // question: by controller? I thought by command-line-tool
	// +optional
	Fingerprint string `json:"fingerprint,omitempty"`
}

// HivelocityMachineType defines the Hivelocity Machine type.
type HivelocityMachineType string

// ResourceLifecycle configures the lifecycle of a resource.
type ResourceLifecycle string


// HivelocitySecretRef defines the name of the Secret and the relevant keys in the secret to access the Hivelocity API.
type HivelocitySecretRef struct {
	// +optional
	// +kubebuilder:default=hivelocity
	Name string `json:"name,omitempty"`

	// +optional
	// +kubebuilder:default=HIVELOCITY_API_KEY
	Key string `json:"key,omitempty"`
}

// PublicNetworkSpec contains specs about public network spec of an Hivelocity server.
type PublicNetworkSpec struct {
	// +optional
	// +kubebuilder:default=true
	EnableIPv4 bool `json:"enableIPv4"`
	// +optional
	// +kubebuilder:default=true
	EnableIPv6 bool `json:"enableIPv6"`
}

// LoadBalancerSpec defines the desired state of the Control Plane Loadbalancer.
type LoadBalancerSpec struct {
	// +optional
	// +kubebuilder:default=true
	Enabled bool `json:"enabled"`

	// +optional
	Name *string `json:"name,omitempty"`

	// Could be round_robin or least_connection. The default value is "round_robin".
	// +optional
	// +kubebuilder:validation:Enum=round_robin;least_connections
	// +kubebuilder:default=round_robin
	Algorithm LoadBalancerAlgorithmType `json:"algorithm,omitempty"`

	// Loadbalancer type
	// +optional
	// +kubebuilder:validation:Enum=lb11;lb21;lb31
	// +kubebuilder:default=lb11
	Type string `json:"type,omitempty"`

	// API Server port. It must be valid ports range (1-65535). If omitted, default value is 6443.
	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:default=6443
	Port int `json:"port,omitempty"`

	// Defines how traffic will be routed from the Load Balancer to your target server.
	// +optional
	ExtraServices []LoadBalancerServiceSpec `json:"extraServices,omitempty"`

	// Region contains the name of the Hivelocity location the load balancer is running.
	Region Region `json:"region,omitempty"`
}

// LoadBalancerServiceSpec defines a Loadbalancer Target.
type LoadBalancerServiceSpec struct {
	// Protocol specifies the supported Loadbalancer Protocol.
	// +kubebuilder:validation:Enum=http;https;tcp
	Protocol string `json:"protocol,omitempty"`

	// ListenPort, i.e. source port, defines the incoming port open on the loadbalancer.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	ListenPort int `json:"listenPort,omitempty"`

	// DestinationPort defines the port on the server.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	DestinationPort int `json:"destinationPort,omitempty"`
}

// LoadBalancerStatus defines the obeserved state of the control plane loadbalancer.
type LoadBalancerStatus struct {
	ID         int                  `json:"id,omitempty"`
	IPv4       string               `json:"ipv4,omitempty"`
	IPv6       string               `json:"ipv6,omitempty"`
	InternalIP string               `json:"internalIP,omitempty"`
	Target     []LoadBalancerTarget `json:"targets,omitempty"`
	Protected  bool                 `json:"protected,omitempty"`
}

// LoadBalancerTarget defines the target of a load balancer.
type LoadBalancerTarget struct {
	Type     LoadBalancerTargetType `json:"type"`
	ServerID int                    `json:"serverID,omitempty"`
	IP       string                 `json:"ip,omitempty"`
}

// HivelocityNetworkSpec defines the desired state of the Hivelocity Private Network.
type HivelocityNetworkSpec struct {
	// Enabled defines whether the network should be enabled or not
	Enabled bool `json:"enabled"`

	// CIDRBlock defines the cidrBlock of the Hivelocity Network. A Subnet is required.
	// +kubebuilder:default="10.0.0.0/16"
	// +optional
	CIDRBlock string `json:"cidrBlock,omitempty"`

	// SubnetCIDRBlock defines the cidrBlock for the subnet of the Hivelocity Network.
	// +kubebuilder:default="10.0.0.0/24"
	// +optional
	SubnetCIDRBlock string `json:"subnetCidrBlock,omitempty"`

	// NetworkZone specifies the Hivelocity network zone of the private network.
	// +kubebuilder:validation:Enum=eu-central;us-east
	// +kubebuilder:default=eu-central
	// +optional
	NetworkZone HivelocityNetworkZone `json:"networkZone,omitempty"`
}

// NetworkStatus defines the observed state of the Hivelocity Private Network.
type NetworkStatus struct {
	ID              int               `json:"id,omitempty"`
	Labels          map[string]string `json:"-"`
	AttachedServers []int             `json:"attachedServers,omitempty"`
}

// Region is a Hivelocity Location
// +kubebuilder:validation:Enum=TODO1;TODO2
type Region string

// HivelocityNetworkZone describes the Network zone.
type HivelocityNetworkZone string

// IsZero returns true if a private Network is set.
func (s *HivelocityNetworkSpec) IsZero() bool {
	if len(s.CIDRBlock) > 0 {
		return false
	}
	if len(s.SubnetCIDRBlock) > 0 {
		return false
	}
	return true
}
