/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type OrderDump struct {
	Info    interface{} `json:"info,omitempty"`
	Total   float32     `json:"total,omitempty"`
	Owner   string      `json:"owner,omitempty"`
	OrderId int32       `json:"orderId,omitempty"`
	Status  string      `json:"status,omitempty"`
}