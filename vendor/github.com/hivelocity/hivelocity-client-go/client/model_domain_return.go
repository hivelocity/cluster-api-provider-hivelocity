/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type DomainReturn struct {
	DomainId  int32       `json:"domainId"`
	DirectsTo string      `json:"directsTo"`
	Summary   interface{} `json:"summary,omitempty"`
	Name      string      `json:"name"`
}