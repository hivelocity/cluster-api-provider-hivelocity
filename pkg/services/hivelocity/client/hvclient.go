// Package hvclient provides the interfaces to communicate with
// the API of Hivelocity.
// We use interfaces to make mocking easier.
package hvclient

import (
	"context"
	"fmt"
	"time"

	apiClient "github.com/hivelocity/hivelocity-client-go/client"
	hv "github.com/hivelocity/hivelocity-client-go/client"
)

const PowerStatusOff = "OFF"
const PowerStatusOn = "ON"

var ErrRateLimitExceeded = fmt.Errorf("Hivelocity API rate limited exceeded")

// Client collects all methods used by the controller in the Hivelocity API.
type Client interface {
	Close()
	PowerOnServer(context.Context, *hv.BareMetalDevice) error
	CreateServer(context.Context, *hv.BareMetalDeviceUpdate) (*hv.BareMetalDevice, error)
}

// Factory is the interface for creating new Client objects.
type Factory interface {
	NewClient(hvAPIKey string) Client
}

type factory struct{}

// NewClient creates new Hivelocity clients.
func (f *factory) NewClient(hvAPIKey string) Client {
	authContext := context.WithValue(context.Background(), apiClient.ContextAPIKey, apiClient.APIKey{
		Key: hvAPIKey,
	})
	apiClient := apiClient.NewAPIClient(apiClient.NewConfiguration())
	return &realClient{
		client:      apiClient,
		authContext: &authContext,
	}
}

type realClient struct {
	client      *apiClient.APIClient
	authContext *context.Context
}

var _ Client = &realClient{}

// Close implements the Close method of the HVClient interface.
func (c *realClient) Close() {}

func (c *realClient) PowerOnServer(context.Context, *hv.BareMetalDevice) error {
	return nil // todo
}

func (c *realClient) CreateServer(context.Context, *hv.BareMetalDeviceUpdate) (*hv.BareMetalDevice, error) {
	return nil, fmt.Errorf("todo")
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
)

type Server struct {
	// todo: returned by findServer()

	ID      int
	Name    string
	Status  ServerStatus
	Created time.Time
	//PublicNet       ServerPublicNet
	//PrivateNet      []ServerPrivateNet
	//ServerType      *ServerType
	//Datacenter      *Datacenter
	IncludedTraffic uint64
	OutgoingTraffic uint64
	IngoingTraffic  uint64
	BackupWindow    string
	RescueEnabled   bool
	Locked          bool
	//ISO             *ISO
	//Image           *Image
	//Protection      ServerProtection
	Labels map[string]string
	//Volumes         []*Volume
	//PrimaryDiskSize int
}

type SSHKey struct {
	// We will only support one ssh-key for the cluster
	// todo
	// There is a second SSHKey struct in api. Maybe one struct is enough?
}
