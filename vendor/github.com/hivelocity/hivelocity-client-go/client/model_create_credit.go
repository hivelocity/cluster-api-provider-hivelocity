/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type CreateCredit struct {
	// The Billing Info ID used to purchase the credits on this account
	BillingInfoId int32 `json:"billingInfoId,omitempty"`
	// The amount of credit associated with the credit ID
	Amount float32 `json:"amount,omitempty"`
}