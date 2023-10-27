<p>Packages:</p>
<ul>
<li>
<a href="#infrastructure.cluster.x-k8s.io%2fv1alpha1">infrastructure.cluster.x-k8s.io/v1alpha1</a>
</li>
</ul>
<h2 id="infrastructure.cluster.x-k8s.io/v1alpha1">infrastructure.cluster.x-k8s.io/v1alpha1</h2>
<p>
<p>Package v1alpha1 contains API Schema definitions for the infrastructure v1alpha1 API group</p>
</p>
Resource Types:
<ul></ul>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.ControllerGeneratedStatus">ControllerGeneratedStatus
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityMachineSpec">HivelocityMachineSpec</a>)
</p>
<p>
<p>ControllerGeneratedStatus contains all status information which is important to persist.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>provisioningState</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.ProvisioningState">
ProvisioningState
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Information tracked by the provisioner.</p>
</td>
</tr>
<tr>
<td>
<code>lastUpdated</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#time-v1-meta">
Kubernetes meta/v1.Time
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Time stamp of last update of status.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityCluster">HivelocityCluster
</h3>
<p>
<p>HivelocityCluster is the Schema for the hivelocityclusters API.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityClusterSpec">
HivelocityClusterSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>controlPlaneEndpoint</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.APIEndpoint
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.</p>
</td>
</tr>
<tr>
<td>
<code>controlPlaneRegion</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.Region">
Region
</a>
</em>
</td>
<td>
<p>ControlPlaneRegion is a Hivelocity Region (LAX2, &hellip;).</p>
</td>
</tr>
<tr>
<td>
<code>hivelocitySecretRef</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocitySecretRef">
HivelocitySecretRef
</a>
</em>
</td>
<td>
<p>HivelocitySecret is a reference to a Kubernetes Secret.</p>
</td>
</tr>
<tr>
<td>
<code>sshKey</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.SSHKey">
SSHKey
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>SSHKey is cluster wide. Valid value is a valid SSH key name.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityClusterStatus">
HivelocityClusterStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityClusterSpec">HivelocityClusterSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityCluster">HivelocityCluster</a>, <a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityClusterTemplateResource">HivelocityClusterTemplateResource</a>)
</p>
<p>
<p>HivelocityClusterSpec defines the desired state of HivelocityCluster.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>controlPlaneEndpoint</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.APIEndpoint
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.</p>
</td>
</tr>
<tr>
<td>
<code>controlPlaneRegion</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.Region">
Region
</a>
</em>
</td>
<td>
<p>ControlPlaneRegion is a Hivelocity Region (LAX2, &hellip;).</p>
</td>
</tr>
<tr>
<td>
<code>hivelocitySecretRef</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocitySecretRef">
HivelocitySecretRef
</a>
</em>
</td>
<td>
<p>HivelocitySecret is a reference to a Kubernetes Secret.</p>
</td>
</tr>
<tr>
<td>
<code>sshKey</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.SSHKey">
SSHKey
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>SSHKey is cluster wide. Valid value is a valid SSH key name.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityClusterStatus">HivelocityClusterStatus
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityCluster">HivelocityCluster</a>)
</p>
<p>
<p>HivelocityClusterStatus defines the observed state of HivelocityCluster.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ready</code><br/>
<em>
bool
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>failureDomains</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.FailureDomains
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>conditions</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.Conditions
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityClusterTemplate">HivelocityClusterTemplate
</h3>
<p>
<p>HivelocityClusterTemplate is the Schema for the hivelocityclustertemplates API.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityClusterTemplateSpec">
HivelocityClusterTemplateSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>template</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityClusterTemplateResource">
HivelocityClusterTemplateResource
</a>
</em>
</td>
<td>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityClusterTemplateResource">HivelocityClusterTemplateResource
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityClusterTemplateSpec">HivelocityClusterTemplateSpec</a>)
</p>
<p>
<p>HivelocityClusterTemplateResource contains spec for HivelocityClusterSpec.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.ObjectMeta
</a>
</em>
</td>
<td>
<em>(Optional)</em>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityClusterSpec">
HivelocityClusterSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>controlPlaneEndpoint</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.APIEndpoint
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.</p>
</td>
</tr>
<tr>
<td>
<code>controlPlaneRegion</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.Region">
Region
</a>
</em>
</td>
<td>
<p>ControlPlaneRegion is a Hivelocity Region (LAX2, &hellip;).</p>
</td>
</tr>
<tr>
<td>
<code>hivelocitySecretRef</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocitySecretRef">
HivelocitySecretRef
</a>
</em>
</td>
<td>
<p>HivelocitySecret is a reference to a Kubernetes Secret.</p>
</td>
</tr>
<tr>
<td>
<code>sshKey</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.SSHKey">
SSHKey
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>SSHKey is cluster wide. Valid value is a valid SSH key name.</p>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityClusterTemplateSpec">HivelocityClusterTemplateSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityClusterTemplate">HivelocityClusterTemplate</a>)
</p>
<p>
<p>HivelocityClusterTemplateSpec defines the desired state of HivelocityClusterTemplate.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>template</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityClusterTemplateResource">
HivelocityClusterTemplateResource
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityDeviceType">HivelocityDeviceType
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityMachineSpec">HivelocityMachineSpec</a>)
</p>
<p>
<p>HivelocityDeviceType defines the Hivelocity device type.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityMachine">HivelocityMachine
</h3>
<p>
<p>HivelocityMachine is the Schema for the hivelocitymachines API.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityMachineSpec">
HivelocityMachineSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>providerID</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ProviderID is the unique identifier as specified by the cloud provider.</p>
</td>
</tr>
<tr>
<td>
<code>type</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityDeviceType">
HivelocityDeviceType
</a>
</em>
</td>
<td>
<p>Type is the Hivelocity Machine Type for this machine.</p>
</td>
</tr>
<tr>
<td>
<code>imageName</code><br/>
<em>
string
</em>
</td>
<td>
<p>ImageName is the reference to the Machine Image from which to create the device.</p>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.ControllerGeneratedStatus">
ControllerGeneratedStatus
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Status contains all status information of the controller. Do not edit these values!</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityMachineStatus">
HivelocityMachineStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityMachineSpec">HivelocityMachineSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityMachine">HivelocityMachine</a>, <a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityMachineTemplateResource">HivelocityMachineTemplateResource</a>)
</p>
<p>
<p>HivelocityMachineSpec defines the desired state of HivelocityMachine.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>providerID</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ProviderID is the unique identifier as specified by the cloud provider.</p>
</td>
</tr>
<tr>
<td>
<code>type</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityDeviceType">
HivelocityDeviceType
</a>
</em>
</td>
<td>
<p>Type is the Hivelocity Machine Type for this machine.</p>
</td>
</tr>
<tr>
<td>
<code>imageName</code><br/>
<em>
string
</em>
</td>
<td>
<p>ImageName is the reference to the Machine Image from which to create the device.</p>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.ControllerGeneratedStatus">
ControllerGeneratedStatus
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Status contains all status information of the controller. Do not edit these values!</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityMachineStatus">HivelocityMachineStatus
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityMachine">HivelocityMachine</a>)
</p>
<p>
<p>HivelocityMachineStatus defines the observed state of HivelocityMachine.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>ready</code><br/>
<em>
bool
</em>
</td>
<td>
<em>(Optional)</em>
<p>Ready is true when the provider resource is ready.</p>
</td>
</tr>
<tr>
<td>
<code>addresses</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
[]Cluster API api/v1beta1.MachineAddress
</a>
</em>
</td>
<td>
<p>Addresses contains the machine&rsquo;s associated addresses.</p>
</td>
</tr>
<tr>
<td>
<code>region</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.Region">
Region
</a>
</em>
</td>
<td>
<p>Region contains the name of the Hivelocity location the device is running.</p>
</td>
</tr>
<tr>
<td>
<code>powerState</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>PowerState is the state of the device for this machine.</p>
</td>
</tr>
<tr>
<td>
<code>failureReason</code><br/>
<em>
<a href="https://pkg.go.dev/sigs.k8s.io/cluster-api@v1.0.0/errors#MachineStatusError">
Cluster API errors.MachineStatusError
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>FailureReason will be set in the event that there is a terminal problem
reconciling the Machine and will contain a succinct value suitable
for machine interpretation.</p>
</td>
</tr>
<tr>
<td>
<code>failureMessage</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>FailureMessage will be set in the event that there is a terminal problem
reconciling the Machine and will contain a more verbose string suitable
for logging and human consumption.</p>
</td>
</tr>
<tr>
<td>
<code>conditions</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.Conditions
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Conditions defines current service state of the HivelocityMachine.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityMachineTemplate">HivelocityMachineTemplate
</h3>
<p>
<p>HivelocityMachineTemplate is the Schema for the hivelocitymachinetemplates API.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityMachineTemplateSpec">
HivelocityMachineTemplateSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>template</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityMachineTemplateResource">
HivelocityMachineTemplateResource
</a>
</em>
</td>
<td>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityMachineTemplateStatus">
HivelocityMachineTemplateStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityMachineTemplateResource">HivelocityMachineTemplateResource
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityMachineTemplateSpec">HivelocityMachineTemplateSpec</a>)
</p>
<p>
<p>HivelocityMachineTemplateResource describes the data needed to create am HivelocityMachine from a template.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.ObjectMeta
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Standard object&rsquo;s metadata.</p>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityMachineSpec">
HivelocityMachineSpec
</a>
</em>
</td>
<td>
<p>Spec is the specification of the desired behavior of the machine.</p>
<br/>
<br/>
<table>
<tr>
<td>
<code>providerID</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ProviderID is the unique identifier as specified by the cloud provider.</p>
</td>
</tr>
<tr>
<td>
<code>type</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityDeviceType">
HivelocityDeviceType
</a>
</em>
</td>
<td>
<p>Type is the Hivelocity Machine Type for this machine.</p>
</td>
</tr>
<tr>
<td>
<code>imageName</code><br/>
<em>
string
</em>
</td>
<td>
<p>ImageName is the reference to the Machine Image from which to create the device.</p>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.ControllerGeneratedStatus">
ControllerGeneratedStatus
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Status contains all status information of the controller. Do not edit these values!</p>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityMachineTemplateSpec">HivelocityMachineTemplateSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityMachineTemplate">HivelocityMachineTemplate</a>)
</p>
<p>
<p>HivelocityMachineTemplateSpec defines the desired state of HivelocityMachineTemplate.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>template</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityMachineTemplateResource">
HivelocityMachineTemplateResource
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityMachineTemplateStatus">HivelocityMachineTemplateStatus
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityMachineTemplate">HivelocityMachineTemplate</a>)
</p>
<p>
<p>HivelocityMachineTemplateStatus defines the observed state of HivelocityMachineTemplate.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>capacity</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#resourcelist-v1-core">
Kubernetes core/v1.ResourceList
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Capacity defines the resource capacity for this machine.
This value is used for autoscaling from zero operations as defined in:
<a href="https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20210310-opt-in-autoscaling-from-zero.md">https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20210310-opt-in-autoscaling-from-zero.md</a></p>
</td>
</tr>
<tr>
<td>
<code>conditions</code><br/>
<em>
<a href="https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api@v1.0.0">
Cluster API api/v1beta1.Conditions
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Conditions defines current service state of the HivelocityMachineTemplate.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityMachineTemplateWebhook">HivelocityMachineTemplateWebhook
</h3>
<p>
<p>HivelocityMachineTemplateWebhook implements a custom validation webhook for HivelocityMachineTemplate.</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityRemediation">HivelocityRemediation
</h3>
<p>
<p>HivelocityRemediation is the Schema for the  hivelocityremediations API.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
<em>(Optional)</em>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityRemediationSpec">
HivelocityRemediationSpec
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<br/>
<br/>
<table>
<tr>
<td>
<code>strategy</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.RemediationStrategy">
RemediationStrategy
</a>
</em>
</td>
<td>
<p>Strategy field defines remediation strategy.</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityRemediationStatus">
HivelocityRemediationStatus
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityRemediationSpec">HivelocityRemediationSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityRemediation">HivelocityRemediation</a>, <a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityRemediationTemplateResource">HivelocityRemediationTemplateResource</a>)
</p>
<p>
<p>HivelocityRemediationSpec defines the desired state of HivelocityRemediation.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>strategy</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.RemediationStrategy">
RemediationStrategy
</a>
</em>
</td>
<td>
<p>Strategy field defines remediation strategy.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityRemediationStatus">HivelocityRemediationStatus
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityRemediation">HivelocityRemediation</a>, <a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityRemediationTemplateStatus">HivelocityRemediationTemplateStatus</a>)
</p>
<p>
<p>HivelocityRemediationStatus defines the observed state of HivelocityRemediation.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>phase</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Phase represents the current phase of machine remediation.
E.g. Pending, Running, Done etc.</p>
</td>
</tr>
<tr>
<td>
<code>retryCount</code><br/>
<em>
int
</em>
</td>
<td>
<em>(Optional)</em>
<p>RetryCount can be used as a counter during the remediation.
Field can hold number of reboots etc.</p>
</td>
</tr>
<tr>
<td>
<code>lastRemediated</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#time-v1-meta">
Kubernetes meta/v1.Time
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>LastRemediated identifies when the host was last remediated</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityRemediationTemplate">HivelocityRemediationTemplate
</h3>
<p>
<p>HivelocityRemediationTemplate is the Schema for the  hivelocityremediationtemplates API.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>metadata</code><br/>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
<em>(Optional)</em>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityRemediationTemplateSpec">
HivelocityRemediationTemplateSpec
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<br/>
<br/>
<table>
<tr>
<td>
<code>template</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityRemediationTemplateResource">
HivelocityRemediationTemplateResource
</a>
</em>
</td>
<td>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityRemediationTemplateStatus">
HivelocityRemediationTemplateStatus
</a>
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityRemediationTemplateResource">HivelocityRemediationTemplateResource
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityRemediationTemplateSpec">HivelocityRemediationTemplateSpec</a>)
</p>
<p>
<p>HivelocityRemediationTemplateResource describes the data needed to create a HivelocityRemediation from a template.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>spec</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityRemediationSpec">
HivelocityRemediationSpec
</a>
</em>
</td>
<td>
<p>Spec is the specification of the desired behavior of the HivelocityRemediation.</p>
<br/>
<br/>
<table>
<tr>
<td>
<code>strategy</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.RemediationStrategy">
RemediationStrategy
</a>
</em>
</td>
<td>
<p>Strategy field defines remediation strategy.</p>
</td>
</tr>
</table>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityRemediationTemplateSpec">HivelocityRemediationTemplateSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityRemediationTemplate">HivelocityRemediationTemplate</a>)
</p>
<p>
<p>HivelocityRemediationTemplateSpec defines the desired state of HivelocityRemediationTemplate.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>template</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityRemediationTemplateResource">
HivelocityRemediationTemplateResource
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityRemediationTemplateStatus">HivelocityRemediationTemplateStatus
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityRemediationTemplate">HivelocityRemediationTemplate</a>)
</p>
<p>
<p>HivelocityRemediationTemplateStatus defines the observed state of HivelocityRemediationTemplate.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityRemediationStatus">
HivelocityRemediationStatus
</a>
</em>
</td>
<td>
<p>HivelocityRemediationStatus defines the observed state of HivelocityRemediation</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.HivelocitySecretRef">HivelocitySecretRef
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityClusterSpec">HivelocityClusterSpec</a>)
</p>
<p>
<p>HivelocitySecretRef defines the name of the Secret and the relevant key in the secret to access the Hivelocity API.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
<tr>
<td>
<code>key</code><br/>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.ProvisioningState">ProvisioningState
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.ControllerGeneratedStatus">ControllerGeneratedStatus</a>)
</p>
<p>
<p>ProvisioningState defines the states the provisioner will report the host has having.</p>
</p>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;associate-device&#34;</p></td>
<td><p>StateAssociateDevice .</p>
</td>
</tr><tr><td><p>&#34;delete&#34;</p></td>
<td><p>StateDeleteDevice .</p>
</td>
</tr><tr><td><p>&#34;delete-deprovision&#34;</p></td>
<td><p>StateDeleteDeviceDeProvision .</p>
</td>
</tr><tr><td><p>&#34;delete-dissociate&#34;</p></td>
<td><p>StateDeleteDeviceDissociate .</p>
</td>
</tr><tr><td><p>&#34;provisioned&#34;</p></td>
<td><p>StateDeviceProvisioned .</p>
</td>
</tr><tr><td><p>&#34;&#34;</p></td>
<td><p>StateNone means the state is unknown.</p>
</td>
</tr><tr><td><p>&#34;provision-device&#34;</p></td>
<td><p>StateProvisionDevice .</p>
</td>
</tr><tr><td><p>&#34;verify-associate&#34;</p></td>
<td><p>StateVerifyAssociate .</p>
</td>
</tr></tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.Region">Region
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityClusterSpec">HivelocityClusterSpec</a>, <a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityMachineStatus">HivelocityMachineStatus</a>)
</p>
<p>
<p>Region is a Hivelocity Location</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.RemediationStrategy">RemediationStrategy
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityRemediationSpec">HivelocityRemediationSpec</a>)
</p>
<p>
<p>RemediationStrategy describes how to remediate machines.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>type</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.RemediationType">
RemediationType
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Type of remediation.</p>
</td>
</tr>
<tr>
<td>
<code>retryLimit</code><br/>
<em>
int
</em>
</td>
<td>
<em>(Optional)</em>
<p>Sets maximum number of remediation retries.</p>
</td>
</tr>
<tr>
<td>
<code>timeout</code><br/>
<em>
<a href="https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1#Duration">
Kubernetes meta/v1.Duration
</a>
</em>
</td>
<td>
<p>Sets the timeout between remediation retries.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.RemediationType">RemediationType
(<code>string</code> alias)</p></h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.RemediationStrategy">RemediationStrategy</a>)
</p>
<p>
<p>RemediationType defines the type of remediation.</p>
</p>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;Reboot&#34;</p></td>
<td><p>RemediationTypeReboot sets RemediationType to Reboot.</p>
</td>
</tr></tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.ResourceLifecycle">ResourceLifecycle
(<code>string</code> alias)</p></h3>
<p>
<p>ResourceLifecycle configures the lifecycle of a resource.</p>
</p>
<table>
<thead>
<tr>
<th>Value</th>
<th>Description</th>
</tr>
</thead>
<tbody><tr><td><p>&#34;owned&#34;</p></td>
<td><p>ResourceLifecycleOwned is the value we use when tagging resources to indicate
that the resource is considered owned and managed by the cluster,
and in particular that the lifecycle is tied to the lifecycle of the cluster.</p>
</td>
</tr><tr><td><p>&#34;shared&#34;</p></td>
<td><p>ResourceLifecycleShared is the value we use when tagging resources to indicate
that the resource is shared between multiple clusters, and should not be destroyed
if the cluster is destroyed.</p>
</td>
</tr></tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.SSHKey">SSHKey
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityClusterSpec">HivelocityClusterSpec</a>)
</p>
<p>
<p>SSHKey defines the SSHKey for Hivelocity.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>name</code><br/>
<em>
string
</em>
</td>
<td>
<p>Name of SSH key.</p>
</td>
</tr>
</tbody>
</table>
<hr/>
