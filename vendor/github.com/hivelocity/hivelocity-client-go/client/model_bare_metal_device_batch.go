/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type BareMetalDeviceBatch struct {
	// List of provisioned devices.
	Devices []BareMetalDevice `json:"devices,omitempty"`
	// Unique ID of the group order.
	OrderGroupId int32 `json:"orderGroupId,omitempty"`
}