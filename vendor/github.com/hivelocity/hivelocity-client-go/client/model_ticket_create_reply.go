/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type TicketCreateReply struct {
	Subject     string        `json:"subject,omitempty"`
	ReplyTo     string        `json:"replyTo,omitempty"`
	ContactId   float32       `json:"contactId,omitempty"`
	Cc          string        `json:"cc,omitempty"`
	Body        string        `json:"body"`
	Type_       float32       `json:"type,omitempty"`
	Encrypted   string        `json:"encrypted,omitempty"`
	Recipient   string        `json:"recipient,omitempty"`
	Headers     string        `json:"headers,omitempty"`
	Attachments []interface{} `json:"attachments,omitempty"`
	Date        float32       `json:"date,omitempty"`
	Hidden      float32       `json:"hidden,omitempty"`
}