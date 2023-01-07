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
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityCluster">HivelocityCluster
</h3>
<p>
<p>HivelocityCluster is the Schema for the hivelocityclusters API</p>
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
<code>foo</code><br/>
<em>
string
</em>
</td>
<td>
<p>Foo is an example field of HivelocityCluster. Edit hivelocitycluster_types.go to remove/update</p>
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
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityCluster">HivelocityCluster</a>)
</p>
<p>
<p>HivelocityClusterSpec defines the desired state of HivelocityCluster</p>
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
<code>foo</code><br/>
<em>
string
</em>
</td>
<td>
<p>Foo is an example field of HivelocityCluster. Edit hivelocitycluster_types.go to remove/update</p>
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
<p>HivelocityClusterStatus defines the observed state of HivelocityCluster</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityClusterTemplate">HivelocityClusterTemplate
</h3>
<p>
<p>HivelocityClusterTemplate is the Schema for the hivelocityclustertemplates API</p>
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
<code>foo</code><br/>
<em>
string
</em>
</td>
<td>
<p>Foo is an example field of HivelocityClusterTemplate. Edit hivelocityclustertemplate_types.go to remove/update</p>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code><br/>
<em>
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityClusterTemplateStatus">
HivelocityClusterTemplateStatus
</a>
</em>
</td>
<td>
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
<p>HivelocityClusterTemplateSpec defines the desired state of HivelocityClusterTemplate</p>
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
<code>foo</code><br/>
<em>
string
</em>
</td>
<td>
<p>Foo is an example field of HivelocityClusterTemplate. Edit hivelocityclustertemplate_types.go to remove/update</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityClusterTemplateStatus">HivelocityClusterTemplateStatus
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityClusterTemplate">HivelocityClusterTemplate</a>)
</p>
<p>
<p>HivelocityClusterTemplateStatus defines the observed state of HivelocityClusterTemplate</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityMachine">HivelocityMachine
</h3>
<p>
<p>HivelocityMachine is the Schema for the hivelocitymachines API</p>
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
<code>foo</code><br/>
<em>
string
</em>
</td>
<td>
<p>Foo is an example field of HivelocityMachine. Edit hivelocitymachine_types.go to remove/update</p>
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
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityMachine">HivelocityMachine</a>)
</p>
<p>
<p>HivelocityMachineSpec defines the desired state of HivelocityMachine</p>
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
<code>foo</code><br/>
<em>
string
</em>
</td>
<td>
<p>Foo is an example field of HivelocityMachine. Edit hivelocitymachine_types.go to remove/update</p>
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
<p>HivelocityMachineStatus defines the observed state of HivelocityMachine</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityMachineTemplate">HivelocityMachineTemplate
</h3>
<p>
<p>HivelocityMachineTemplate is the Schema for the hivelocitymachinetemplates API</p>
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
<code>foo</code><br/>
<em>
string
</em>
</td>
<td>
<p>Foo is an example field of HivelocityMachineTemplate. Edit hivelocitymachinetemplate_types.go to remove/update</p>
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
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityMachineTemplateSpec">HivelocityMachineTemplateSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityMachineTemplate">HivelocityMachineTemplate</a>)
</p>
<p>
<p>HivelocityMachineTemplateSpec defines the desired state of HivelocityMachineTemplate</p>
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
<code>foo</code><br/>
<em>
string
</em>
</td>
<td>
<p>Foo is an example field of HivelocityMachineTemplate. Edit hivelocitymachinetemplate_types.go to remove/update</p>
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
<p>HivelocityMachineTemplateStatus defines the observed state of HivelocityMachineTemplate</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityRemediation">HivelocityRemediation
</h3>
<p>
<p>HivelocityRemediation is the Schema for the hivelocityremediations API</p>
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
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityRemediationSpec">
HivelocityRemediationSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>foo</code><br/>
<em>
string
</em>
</td>
<td>
<p>Foo is an example field of HivelocityRemediation. Edit hivelocityremediation_types.go to remove/update</p>
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
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityRemediationSpec">HivelocityRemediationSpec
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityRemediation">HivelocityRemediation</a>)
</p>
<p>
<p>HivelocityRemediationSpec defines the desired state of HivelocityRemediation</p>
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
<code>foo</code><br/>
<em>
string
</em>
</td>
<td>
<p>Foo is an example field of HivelocityRemediation. Edit hivelocityremediation_types.go to remove/update</p>
</td>
</tr>
</tbody>
</table>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityRemediationStatus">HivelocityRemediationStatus
</h3>
<p>
(<em>Appears on:</em><a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityRemediation">HivelocityRemediation</a>)
</p>
<p>
<p>HivelocityRemediationStatus defines the observed state of HivelocityRemediation</p>
</p>
<h3 id="infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityRemediationTemplate">HivelocityRemediationTemplate
</h3>
<p>
<p>HivelocityRemediationTemplate is the Schema for the hivelocityremediationtemplates API</p>
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
<a href="#infrastructure.cluster.x-k8s.io/v1alpha1.HivelocityRemediationTemplateSpec">
HivelocityRemediationTemplateSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>foo</code><br/>
<em>
string
</em>
</td>
<td>
<p>Foo is an example field of HivelocityRemediationTemplate. Edit hivelocityremediationtemplate_types.go to remove/update</p>
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
<p>HivelocityRemediationTemplateSpec defines the desired state of HivelocityRemediationTemplate</p>
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
<code>foo</code><br/>
<em>
string
</em>
</td>
<td>
<p>Foo is an example field of HivelocityRemediationTemplate. Edit hivelocityremediationtemplate_types.go to remove/update</p>
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
<p>HivelocityRemediationTemplateStatus defines the observed state of HivelocityRemediationTemplate</p>
</p>
<hr/>
