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

// Package csr contains functions to validate certificate signing requests.
package csr

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"fmt"
	"reflect"
	"strings"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// nodesPrefix defines the prefix name for a node.
const nodesPrefix = "system:node:"

// nodesGroup defines the group name for a node.
const nodesGroup = "system:nodes"

// ValidateKubeletCSR validates a CSR.
func ValidateKubeletCSR(csr *x509.CertificateRequest, machineName string, addresses []clusterv1.MachineAddress) error {
	// check signature and exist quickly
	if err := csr.CheckSignature(); err != nil {
		return fmt.Errorf("failed to check signature of x509 certificate: %w", err)
	}

	var multierr error

	// validate subject
	username := nodesPrefix + machineName

	subjectExpected := pkix.Name{
		CommonName:   username,
		Organization: []string{nodesGroup},
		Names: []pkix.AttributeTypeAndValue{
			{Type: asn1.ObjectIdentifier{2, 5, 4, 10}, Value: nodesGroup},
			{Type: asn1.ObjectIdentifier{2, 5, 4, 3}, Value: username},
		},
	}
	if !reflect.DeepEqual(subjectExpected, csr.Subject) {
		multierr = errors.Join(fmt.Errorf("unexpected subject actual=%+#v, expected=%+#v", csr.Subject, subjectExpected))
	}

	// check for DNS Names
	if len(csr.EmailAddresses) > 0 {
		multierr = errors.Join(fmt.Errorf("email addresses are not allow on the request: %v", csr.EmailAddresses))
	}

	// allow only certain DNS names
	for _, name := range csr.DNSNames {
		if name != machineName {
			multierr = errors.Join(fmt.Errorf("the DNS name %q is not allowed", name))
		}
	}

	// allow only certain IP addresses
	allowedIPAddresses := make(map[string]struct{})
	for _, address := range addresses {
		switch address.Type {
		case clusterv1.MachineInternalIP, clusterv1.MachineExternalIP:
			allowedIPAddresses[strings.Split(address.Address, "/")[0]] = struct{}{}
		}
	}

	for _, ip := range csr.IPAddresses {
		if _, ok := allowedIPAddresses[ip.String()]; !ok {
			multierr = errors.Join(fmt.Errorf("the IP address %q is not allowed", ip.String()))
		}
	}

	return multierr
}
