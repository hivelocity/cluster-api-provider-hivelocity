/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type CreateEnterpriseControlledClient struct {
	// The company to be asociated with the client account
	Company string `json:"company"`
	// The email to be associated with the client account
	Email string `json:"email"`
	// The client account password to be used
	Password string `json:"password,omitempty"`
}
