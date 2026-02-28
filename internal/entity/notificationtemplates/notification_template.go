package notificationtemplates

import "time"

type NotificationTemplate struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Subject      string                 `json:"subject,omitempty"`
	Body         string                 `json:"body,omitempty"`
	PayloadSchema map[string]interface{} `json:"payload_schema,omitempty"`
	Channel      string                 `json:"channel"`
	TemplateType string                 `json:"template_type,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    *time.Time             `json:"updated_at,omitempty"`
	DeletedAt    *time.Time             `json:"deleted_at,omitempty"`
}
