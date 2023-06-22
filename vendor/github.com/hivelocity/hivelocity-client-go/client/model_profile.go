/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type Profile struct {
	Address  interface{} `json:"address,omitempty"`
	Email    string      `json:"email,omitempty"`
	Zip      interface{} `json:"zip,omitempty"`
	MetaData interface{} `json:"metaData,omitempty"`
	Fax      interface{} `json:"fax,omitempty"`
	Id       int32       `json:"id,omitempty"`
	First    string      `json:"first,omitempty"`
	FullName interface{} `json:"fullName,omitempty"`
	City     interface{} `json:"city,omitempty"`
	Country  interface{} `json:"country,omitempty"`
	Last     string      `json:"last,omitempty"`
	Phone    string      `json:"phone,omitempty"`
	State    interface{} `json:"state,omitempty"`
	Company  interface{} `json:"company,omitempty"`
	Login    string      `json:"login,omitempty"`
	IsClient bool        `json:"isClient,omitempty"`
	Created  interface{} `json:"created,omitempty"`
}
