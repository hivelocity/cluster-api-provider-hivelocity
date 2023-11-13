/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type ARecord struct {
	Name      string   `json:"name"`
	Ttl       int32    `json:"ttl"`
	Addresses []string `json:"addresses,omitempty"`
}