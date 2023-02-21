// Package hvclient provides the interfaces to communicate with
// the API of Hivelocity.
// We use interfaces to make mocking easier.
package hvclient

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	hv "github.com/hivelocity/hivelocity-client-go/client"
)

// PowerStatusOff is "OFF".
const PowerStatusOff = "OFF"

// PowerStatusOn is "ON".
const PowerStatusOn = "ON"

// Client collects all methods used by the controller in the Hivelocity API.
type Client interface {
	Close()
	PowerOnServer(ctx context.Context, deviceID int32) error
	CreateServer(ctx context.Context, deviceID int32, opts hv.BareMetalDeviceUpdate) (hv.BareMetalDevice, error)
	ListServers(context.Context) ([]*hv.BareMetalDevice, error)
	ShutdownServer(ctx context.Context, deviceID int32) error
	DeleteServer(ctx context.Context, deviceID int32) error
	ListImages(ctx context.Context, productID int32) ([]string, error)
	ListSSHKeys(context.Context) ([]hv.SshKeyResponse, error)
}

// Factory is the interface for creating new Client objects.
type Factory interface {
	NewClient(hvAPIKey string) Client
}

// HivelocityFactory implements the Factory interface.
type HivelocityFactory struct{}

// ErrDeviceNotFound gets returned if no matching device was found.
var ErrDeviceNotFound = fmt.Errorf("device was not found")

var _ Factory = &HivelocityFactory{}

// NewClient creates new Hivelocity clients.
func (f *HivelocityFactory) NewClient(hvAPIKey string) Client {
	config := hv.NewConfiguration()
	config.AddDefaultHeader("X-API-KEY", hvAPIKey)
	apiClient := hv.NewAPIClient(config)
	return &realClient{
		client: apiClient,
	}
}

type realClient struct {
	client *hv.APIClient
}

var _ Client = &realClient{}

// Close implements the Close method of the HVClient interface.
func (c *realClient) Close() {}

func (c *realClient) PowerOnServer(ctx context.Context, deviceID int32) error {
	return nil // todo
}

func (c *realClient) CreateServer(ctx context.Context, deviceID int32, opts hv.BareMetalDeviceUpdate) (hv.BareMetalDevice, error) {
	// https://developers.hivelocity.net/reference/put_bare_metal_device_id_resource
	device, _, err := c.client.BareMetalDevicesApi.PutBareMetalDeviceIdResource(ctx, deviceID, opts, nil)
	return device, err
}

func (c *realClient) ListServers(ctx context.Context) ([]*hv.BareMetalDevice, error) {
	servers, _, err := c.client.BareMetalDevicesApi.GetBareMetalDeviceResource(ctx, nil)
	ret := make([]*hv.BareMetalDevice, 0, len(servers))
	for i := range servers {
		ret = append(ret, &servers[i])
	}
	return ret, err
}

func (c *realClient) DeleteServer(ctx context.Context, deviceID int32) error {
	return fmt.Errorf("todo DeleteServer")
}

func (c *realClient) ShutdownServer(ctx context.Context, deviceID int32) error {
	return fmt.Errorf("todo ShutdownServer")
}

func (c *realClient) ListImages(ctx context.Context, productID int32) ([]string, error) {
	// https://developers.hivelocity.net/reference/get_product_operating_systems_resource
	opts, _, err := c.client.ProductApi.GetProductOperatingSystemsResource(ctx, productID, nil)
	ret := make([]string, 0, len(opts))
	if err != nil {
		return []string{}, err
	}
	for i := range opts {
		ret = append(ret, opts[i].Name)
	}
	return ret, nil
}

func (c *realClient) ListSSHKeys(ctx context.Context) ([]hv.SshKeyResponse, error) {
	// https://developers.hivelocity.net/reference/get_ssh_key_resource
	sshKeys, _, err := c.client.SshKeyApi.GetSshKeyResource(ctx, nil)
	return sshKeys, err
}

// ServerStatus specifies a server's status.
type ServerStatus string

const (
	// ServerStatusInitializing is the status when a server is initializing.
	ServerStatusInitializing ServerStatus = "initializing" // TODO AFAIK HV does not provide these detailed infos

	// ServerStatusOff is the status when a server is off.
	ServerStatusOff ServerStatus = "off"

	// ServerStatusRunning is the status when a server is running.
	ServerStatusRunning ServerStatus = "running"

	// ServerStatusStarting is the status when a server is being started.
	ServerStatusStarting ServerStatus = "starting"

	// ServerStatusStopping is the status when a server is being stopped.
	ServerStatusStopping ServerStatus = "stopping"

	// ServerStatusMigrating is the status when a server is being migrated.
	ServerStatusMigrating ServerStatus = "migrating"

	// ServerStatusRebuilding is the status when a server is being rebuilt.
	ServerStatusRebuilding ServerStatus = "rebuilding"

	// ServerStatusDeleting is the status when a server is being deleted.
	ServerStatusDeleting ServerStatus = "deleting"

	// ServerStatusUnknown is the status when a server's state is unknown.
	ServerStatusUnknown ServerStatus = "unknown"

	// TagKeyMachineName is the prefix for HV tags for machine names.
	TagKeyMachineName = "caphv-machine-name"

	// TagKeyClusterName is the prefix for HV tags for cluster names.
	TagKeyClusterName = "caphv-cluster-name"

	// TagKeyInstanceType is the prefix for HV tags for instance types.
	TagKeyInstanceType = "caphv-instance-type"
)

// GetMachineTag create tag for HV API. Example: "mymachine" --> "caphv-machine-name=mymachine".
func GetMachineTag(machineName string) string {
	return TagKeyMachineName + "=" + machineName
}

// GetClusterTag create tag for HV API. Example: "mycluster" --> "caphv-cluster-name=mycluster".
func GetClusterTag(clusterName string) string {
	return TagKeyClusterName + "=" + clusterName
}

// IsRateLimitExceededError returns true, if the Hivelocity rate limit was reached.
func IsRateLimitExceededError(err error) bool {
	var swaggerErr hv.GenericSwaggerError
	if !errors.As(err, &swaggerErr) {
		return false
	}
	if strings.HasPrefix(swaggerErr.Error(), fmt.Sprint(http.StatusTooManyRequests)) {
		return true
	}
	return false
}
