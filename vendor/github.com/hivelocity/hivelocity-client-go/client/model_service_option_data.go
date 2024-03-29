/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type ServiceOptionData struct {
	ServiceOptionId int32  `json:"serviceOptionId,omitempty"`
	Name            string `json:"name,omitempty"`
	OptionId        int32  `json:"optionId,omitempty"`
	UpgradeName     string `json:"upgradeName,omitempty"`
	GroupName       string `json:"groupName,omitempty"`
	InvoiceHide     string `json:"invoiceHide,omitempty"`
}
