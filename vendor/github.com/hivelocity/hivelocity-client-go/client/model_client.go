/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type Client struct {
	Zip      string      `json:"zip,omitempty"`
	Last     string      `json:"last,omitempty"`
	Id       float32     `json:"id,omitempty"`
	Email    string      `json:"email,omitempty"`
	Login    string      `json:"login,omitempty"`
	FullName string      `json:"fullName,omitempty"`
	Company  string      `json:"company,omitempty"`
	State    string      `json:"state,omitempty"`
	Country  string      `json:"country,omitempty"`
	MetaData interface{} `json:"metaData,omitempty"`
	IsClient bool        `json:"isClient,omitempty"`
	First    string      `json:"first,omitempty"`
	City     string      `json:"city,omitempty"`
}
