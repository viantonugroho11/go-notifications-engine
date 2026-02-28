package dto

import (
	tplEntity "go-boilerplate-clean/internal/entity/notificationtemplates"
)

type CreateNotificationTemplateRequest struct {
	Name          string                 `json:"name"`
	Subject       string                 `json:"subject,omitempty"`
	Body          string                 `json:"body,omitempty"`
	PayloadSchema map[string]interface{} `json:"payload_schema,omitempty"`
	Channel       string                 `json:"channel"`
	TemplateType  string                 `json:"template_type,omitempty"`
}

type UpdateNotificationTemplateRequest struct {
	Name          string                 `json:"name"`
	Subject       string                 `json:"subject,omitempty"`
	Body          string                 `json:"body,omitempty"`
	PayloadSchema map[string]interface{} `json:"payload_schema,omitempty"`
	Channel       string                 `json:"channel"`
	TemplateType  string                 `json:"template_type,omitempty"`
}

func (r *CreateNotificationTemplateRequest) ToEntity() tplEntity.NotificationTemplate {
	return tplEntity.NotificationTemplate{
		Name:          r.Name,
		Subject:       r.Subject,
		Body:          r.Body,
		PayloadSchema: r.PayloadSchema,
		Channel:       r.Channel,
		TemplateType:  r.TemplateType,
	}
}
