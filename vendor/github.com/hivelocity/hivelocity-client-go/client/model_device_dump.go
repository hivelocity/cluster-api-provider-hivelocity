/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type DeviceDump struct {
	// The unique ID of the device.
	DeviceId int32 `json:"deviceId,omitempty"`
	// User given custom name.
	Name string `json:"name"`
	// active|inactive
	Status string `json:"status,omitempty"`
	// Generic description of device. Usually type and rack unit size.
	DeviceType string `json:"deviceType,omitempty"`
	// Generic group
	DeviceTypeGroup string `json:"deviceTypeGroup,omitempty"`
	// ON|OFF
	PowerStatus interface{} `json:"powerStatus,omitempty"`
	// True if device has active cancellation request.
	HasCancellation bool `json:"hasCancellation,omitempty"`
	// True if device enrolled in managed services.
	IsManaged bool `json:"isManaged,omitempty"`
	// True if device currently reloading.
	IsReload bool `json:"isReload,omitempty"`
	// # of passing device monitors
	MonitorsUp int32 `json:"monitorsUp,omitempty"`
	// Total # device monitors
	MonitorsTotal int32 `json:"monitorsTotal,omitempty"`
	// # of managed service alerts.
	ManagedAlertsTotal int32 `json:"managedAlertsTotal,omitempty"`
	// Device Ports (Network Interfaces).
	Ports []interface{} `json:"ports,omitempty"`
	// a fqdn for the device. for example: `example.hivelocity.net`.
	Hostname string `json:"hostname,omitempty"`
	// True if device is accessible over IPMI by customer.
	IpmiEnabled bool `json:"ipmiEnabled,omitempty"`
	// List containing key/values of device info based on tag order.
	DisplayedTags []interface{} `json:"displayedTags,omitempty"`
	// List of all user set device tags.
	Tags []string `json:"tags,omitempty"`
	// Detailed information on the device location.
	Location interface{} `json:"location,omitempty"`
	// Network Automation status for device.
	NetworkAutomation interface{} `json:"networkAutomation,omitempty"`
	// The first assigned public IP for accessing this device.
	PrimaryIp string `json:"primaryIp,omitempty"`
	// IP address for IPMI connection. Requires you to whitelist your current IP or be on IPMI VPN.
	IpmiAddress     interface{} `json:"ipmiAddress,omitempty"`
	ServiceMonitors []string    `json:"serviceMonitors,omitempty"`
	// If set, detailed info on this device's billing method. Otherwise null. When null the accounts default billing info is used for payments.
	BillingInfo interface{} `json:"billingInfo,omitempty"`
	// The unique ID of the associated service.
	ServicePlan int32 `json:"servicePlan,omitempty"`
	// The unique ID of the last invoice for this device.
	LastInvoiceId int32 `json:"lastInvoiceId,omitempty"`
	// True if instant server.
	SelfProvisioning bool `json:"selfProvisioning,omitempty"`
	// Additional metadata.
	Metadata interface{} `json:"metadata,omitempty"`
	// BUILDING|IPMI_READY|PROVISIONABLE|RESERVED|WAIT_FOR_PXE|PROVISION_STARTED|PROVISION_WAIT_FOR_ADDONS|PROVISION_FINISHED|WAIT_TO_COMPLETE_ORDER|WAIT_TO_ASSIGN_SERVICE|WAIT_FOR_HARDWARE_SCAN|IN_USE|RELOADING|DEVICE_READY_TO_TEST|DEVICE_READY_TO_WIPE|DEVICE_READY_TO_UPGRADE_FIRMWARE|FAILED|CLEANUP_MOVE_TO_FAILED|IN_REVIEW
	SpsStatus string `json:"spsStatus,omitempty"`
}