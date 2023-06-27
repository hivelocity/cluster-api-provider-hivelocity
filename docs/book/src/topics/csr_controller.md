# CSR Controller

> In public key infrastructure (PKI) systems, a certificate signing request (also CSR or certification request) is a message sent
> from an applicant to a certificate authority of the public key infrastructure in order to apply for a digital identity
> certificate. It usually contains the public key for which the certificate should be issued, identifying information (such as a
> domain name) and a proof of authenticity including integrity protection (e.g., a digital signature).

Source: [Wikipedia](https://en.wikipedia.org/wiki/Certificate_signing_request)


Quoting [Enabling signed kubelet serving certificates](https://kubernetes.io/docs/tasks/administer-cluster/kubeadm/kubeadm-certs/#kubelet-serving-certs)

> By default the kubelet serving certificate deployed by kubeadm is self-signed. This means a connection from external services
> like the metrics-server to a kubelet cannot be secured with TLS.
> ... One known limitation is that the CSRs (Certificate Signing Requests) for these certificates cannot be automatically
> approved by the default signer in the kube-controller-manager - kubernetes.io/kubelet-serving. This will require action from
> the user or a third party controller.

The CAPHV CSR Controller signs kubelet-serving certs, since this is not done by Kubernetes or kubeadm up to now.

Related: [List of Kubernetes Signers](https://kubernetes.io/docs/reference/access-authn-authz/certificate-signing-requests/#kubernetes-signers)

The good news for you: You don't need to do anything. It is enabled by default.

The CAPHV CSR controller will automatically sign the CSRs of the kubelets.

Alternative solutions would be:

* Use a tool like [postfinance/kubelet-csr-approver](https://github.com/postfinance/kubelet-csr-approver)
* Use self-signed certs and access the kubelet with insecure TLS. For example metrics-server.
* Approve CSR by hand (kubectl)
