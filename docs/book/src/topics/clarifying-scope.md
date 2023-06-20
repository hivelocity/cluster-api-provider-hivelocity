# Clarifying Scope

Managing a production-grade Kubernetes system requires a seasoned team of experts.

The Cluster API Provider Hivelocity (CAPHV) plays a specific, limited role within this intricate ecosystem. It deals with the low-level lifecycle management of machines and infrastructure. While it can support the creation of production-ready Kubernetes clusters, it does not define the specific contents within the cluster. This important task is something that needs to be handled independently.

Please bear in mind, the standalone utility of this software is not intended to be production-ready!

Here are several aspects that CAPHV will not handle for you:

- production-ready node images
- secured kubeadm configuration
- incorporation of cluster add-ons, such as CNI or metrics-server
- testing & update procedures
- backup procedures
- monitoring strategies
- alerting systems
- identity and Access Management (IAM)

The Cluster API Provider Hivelocity simply equips you with the capacity to manage a Kubernetes cluster using the Cluster API (CAPI). The security measures and defining characteristics of your Kubernetes cluster have to be addressed by you, or expert professionals should be consulted for the same.

Please note that ready-to-use Kubernetes configurations, production-ready node images, kubeadm configuration, cluster add-ons like CNI and similar services need to be separately prepared or acquired to ensure a comprehensive and secure Kubernetes deployment.
