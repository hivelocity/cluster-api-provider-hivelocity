/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type TicketSeachResult struct {
	PerPage  int32       `json:"perPage,omitempty"`
	Page     int32       `json:"page,omitempty"`
	Pages    int32       `json:"pages,omitempty"`
	Total    int32       `json:"total,omitempty"`
	NextPage int32       `json:"nextPage,omitempty"`
	HasNext  bool        `json:"hasNext,omitempty"`
	HasPrev  bool        `json:"hasPrev,omitempty"`
	PrevPage int32       `json:"prevPage,omitempty"`
	Items    *TicketPost `json:"items,omitempty"`
}
