package notificationtemplates

import "time"

type NotificationTemplate struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	Subject       string         `json:"subject,omitempty"`
	Body          string         `json:"body,omitempty"`
	PayloadSchema map[string]any `json:"payload_schema,omitempty"`
	Channel       string         `json:"channel"`
	TemplateType  TemplateType   `json:"template_type,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     *time.Time     `json:"updated_at,omitempty"`
	DeletedAt     *time.Time     `json:"deleted_at,omitempty"`
}

type TemplateType string

const (
	TemplateTypePromo         TemplateType = "promo"
	TemplateTypeTransactional TemplateType = "transactional"
	TemplateTypeSystem        TemplateType = "system"
	TemplateTypeOther         TemplateType = "other"
)

func (t TemplateType) String() string {
	return string(t)
}

// NotificationTemplateListParam dipakai untuk filter dan pagination List notification template.
type NotificationTemplateListParam struct {
	IDs           []string
	Name          string
	Channel       string
	TemplateTypes []string
	Limit, Offset int
}
