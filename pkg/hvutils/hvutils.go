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

// Package hvutils implements helper functions for the HV API.
package hvutils

import (
	"errors"
	"fmt"
	"strings"

	hvclient "github.com/hivelocity/cluster-api-provider-hivelocity/pkg/services/hivelocity/client"
	hv "github.com/hivelocity/hivelocity-client-go/client"
	"golang.org/x/exp/slices"
)

var errMultipleServerFound = fmt.Errorf(
	"multiple servers found while trying to find a single server",
)

// FindServerByTags returns the server with the given clusterTag and machineTag.
// Returns nil if no server was found.
func FindServerByTags(
	clusterTag string,
	machineTag string,
	servers []*hv.BareMetalDevice,
) (*hv.BareMetalDevice, error) {
	var server *hv.BareMetalDevice
	found := 0
	for i := range servers {
		if slices.Contains(servers[i].Tags, clusterTag) &&
			slices.Contains(servers[i].Tags, machineTag) {
			server = servers[i]
			found++
		}
	}
	if found > 1 {
		return nil, fmt.Errorf("found %v servers with tags %s and %s. Expected one: %w",
			found, clusterTag, machineTag, errMultipleServerFound)
	} else if found == 0 {
		return nil, nil
	}
	return server, nil
}

// FindUnusedServer returns an unused server. Returns nil if no server was found.
func FindUnusedServer(servers []*hv.BareMetalDevice, clusterName string, instanceType string) (*hv.BareMetalDevice, error) {
	for i := range servers {
		server := servers[i]
		it, err := GetInstanceType(server)
		if err != nil {
			return nil, fmt.Errorf("[FindUnusedServer] GetInstanceType() failed: %w", err)
		}
		if it != instanceType {
			continue
		}
		if ServerHasTagKey(server, hvclient.TagKeyMachineName) {
			continue
		}
		cn, err := ServerGetTagValue(server, hvclient.TagKeyClusterName)
		if errors.Is(err, ErrTooManyTagsFound) {
			continue
		}
		if errors.Is(err, ErrNoMatchingTagFound) {
			// this could lead to a race-condition, if two controllers of two clusters
			// try to fetch an unused server.
			// TODO: re-check after N seconds if there is a second tag from a second controller.
			return server, nil
		}
		if err != nil {
			return nil, err
		}
		if cn != clusterName {
			continue
		}
		return server, nil
	}
	return nil, nil
}

// ServerHasTagKey returns true if the server has the tagKey set.
// Example: Your can check if a machine has already a name by using tagKey="machine-name".
func ServerHasTagKey(server *hv.BareMetalDevice, tagKey string) bool {
	prefix := tagKey + "="
	for i := range server.Tags {
		if strings.HasPrefix(server.Tags[i], prefix) {
			return true
		}
	}
	return false
}

// ErrTooManyTagsFound gets returned, if there are multiple tags with the same key,
// and the key should be unique.
var ErrTooManyTagsFound = fmt.Errorf("too many tags found")

// ErrNoMatchingTagFound gets returned, if no matching tag was found.
var ErrNoMatchingTagFound = fmt.Errorf("no matching tag found")

// ServerGetTagValue returns the value of a TagKey of a server.
// Example: If a server has the tag "foo=bar", then ServerGetTagValue
// will return "bar".
// If there is no such tag, or there are two tags, then an error gets returned.
func ServerGetTagValue(server *hv.BareMetalDevice, tagKey string) (string, error) {
	prefix := tagKey + "="
	found := 0
	value := ""
	for i := range server.Tags {
		if !strings.HasPrefix(server.Tags[i], prefix) {
			continue
		}
		if found > 0 {
			return "", fmt.Errorf("[ServerGetTagValue] device %q, tagKey %q: %w",
				server.Hostname, tagKey, ErrTooManyTagsFound)
		}
		found++
		value = server.Tags[i][len(prefix):]
	}
	if found == 0 {
		return "", fmt.Errorf("[ServerGetTagValue] device %q, tagKey %q: %w",
			server.Hostname, tagKey, ErrNoMatchingTagFound)
	}
	return value, nil
}

// GetInstanceType returns the instance-type of this BareMetalDevice.
func GetInstanceType(server *hv.BareMetalDevice) (string, error) {
	instanceType, err := ServerGetTagValue(server, hvclient.TagKeyInstanceType)
	if err != nil {
		return "", fmt.Errorf("[GetInstanceType] ServerGetTagValue() failed: %w", err)
	}
	return instanceType, nil
}
