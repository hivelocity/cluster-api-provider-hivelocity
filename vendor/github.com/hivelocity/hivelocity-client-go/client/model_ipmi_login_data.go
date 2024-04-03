/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type IpmiLoginData struct {
	Drivertype string `json:"drivertype,omitempty"`
	// Username for IPMI console.
	Username string `json:"username,omitempty"`
	// IP for IPMI access. Requires your current IP to be whitelisted or the IPMI VPN.
	Host string `json:"host,omitempty"`
	// Password for IPMI console.
	Password string `json:"password,omitempty"`
}