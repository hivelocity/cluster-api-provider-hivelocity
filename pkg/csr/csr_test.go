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

package csr_test

import (
	"crypto/x509"
	"encoding/pem"

	"github.com/hivelocity/cluster-api-provider-hivelocity/pkg/csr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

var _ = Describe("Validate Kubelet CSR", func() {
	var cr *x509.CertificateRequest
	var name string
	var addresses []clusterv1.MachineAddress
	BeforeEach(func() {
		name = "hvtest-guettli-md-0-dsls9"
		addresses = []clusterv1.MachineAddress{
			{
				Type:    clusterv1.MachineExternalIP,
				Address: "66.206.8.194",
			},
		}
		var csrString = `-----BEGIN CERTIFICATE REQUEST-----
MIIBQDCB5gIBADBHMRUwEwYDVQQKEwxzeXN0ZW06bm9kZXMxLjAsBgNVBAMTJXN5
c3RlbTpub2RlOmh2dGVzdC1ndWV0dGxpLW1kLTAtZHNsczkwWTATBgcqhkjOPQIB
BggqhkjOPQMBBwNCAATWsiJI5U2/Y4qWr4PkOC7A92R5g9lslnRKfyTnPR/Cm7Ub
338yswLj9yyc7y0jpHCxK2nUURkvbXFZxOoY0zt2oD0wOwYJKoZIhvcNAQkOMS4w
LDAqBgNVHREEIzAhghlodnRlc3QtZ3VldHRsaS1tZC0wLWRzbHM5hwRCzgjCMAoG
CCqGSM49BAMCA0kAMEYCIQC9BlA2uA5xULv07z/tGxfjWrYbEC0dfwSmjyJd5Aa7
XwIhAMOhY6EyOQX3j356MVa0b2ixNKy9EDtMZbAWVXrClnRl
-----END CERTIFICATE REQUEST-----`
		block, _ := pem.Decode([]byte(csrString))
		var err error
		cr, err = x509.ParseCertificateRequest(block.Bytes)
		Expect(err).To(BeNil())
	})

	It("should not fail", func() {
		Expect(csr.ValidateKubeletCSR(cr, name, addresses)).To(Succeed())
	})
})
