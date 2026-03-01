package model

import (
	"encoding/json"
	notificationtemplates "go-boilerplate-clean/internal/entity/notificationtemplates"
	"time"
)

type NotificationTemplate struct {
	ID            string     `gorm:"column:id;primaryKey"`
	Name          string     `gorm:"column:name;not null"`
	Subject       string     `gorm:"column:subject"`
	Body          string     `gorm:"column:body;type:text"`
	PayloadSchema []byte     `gorm:"column:payload_schema;type:jsonb"`
	Channel       string     `gorm:"column:channel;not null"`
	TemplateType  string     `gorm:"column:template_type"`
	CreatedAt     time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt     *time.Time `gorm:"column:updated_at"`
	DeletedAt     *time.Time `gorm:"column:deleted_at"`
}

func (NotificationTemplate) TableName() string {
	return "notification_templates"
}


func (n NotificationTemplate) ToEntity() notificationtemplates.NotificationTemplate {
	var payloadSchema map[string]interface{}
	if len(n.PayloadSchema) > 0 {
		_ = json.Unmarshal(n.PayloadSchema, &payloadSchema)
	}
	return notificationtemplates.NotificationTemplate{
		ID: n.ID,
		Name: n.Name,
		Subject: n.Subject,
		Body: n.Body,
		PayloadSchema: payloadSchema,
		Channel: n.Channel,
		TemplateType: n.TemplateType,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
		DeletedAt: n.DeletedAt,
	}
}

func ToDBNotificationTemplate(t notificationtemplates.NotificationTemplate) NotificationTemplate {
	schemaJSON, _ := json.Marshal(t.PayloadSchema)
	return NotificationTemplate{
		ID: t.ID,
		Name: t.Name,
		Subject: t.Subject,
		Body: t.Body,
		PayloadSchema: schemaJSON,
		Channel: t.Channel,
		TemplateType: t.TemplateType,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
		DeletedAt: t.DeletedAt,
	}
}