/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type ClientCreateDump struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Fax      string `json:"fax,omitempty"`
	Id       int32  `json:"id,omitempty"`
	Phone    string `json:"phone,omitempty"`
	City     string `json:"city,omitempty"`
	Active   bool   `json:"active,omitempty"`
	Last     string `json:"last,omitempty"`
	State    string `json:"state,omitempty"`
	Address  string `json:"address,omitempty"`
	Company  string `json:"company,omitempty"`
	Created  int32  `json:"created,omitempty"`
	Zip      string `json:"zip,omitempty"`
	First    string `json:"first,omitempty"`
	Country  string `json:"country,omitempty"`
}
