package model

import (
	"encoding/json"
	notificationtemplates "go-boilerplate-clean/internal/entity/notificationtemplates"
	"time"

	"gorm.io/gorm"
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
		TemplateType: notificationtemplates.TemplateType(n.TemplateType),
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
		TemplateType: t.TemplateType.String(),
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
		DeletedAt: t.DeletedAt,
	}
}

// ApplyListParam menerapkan filter dan pagination dari param ke query. Jika param nil, query tidak diubah.
func ApplyListParam(query *gorm.DB, param *notificationtemplates.NotificationTemplateListParam) *gorm.DB {
	if param == nil {
		return query
	}
	if len(param.IDs) > 0 {
		query = query.Where("id IN (?)", param.IDs)
	}
	if param.Name != "" {
		query = query.Where("name = ?", param.Name)
	}
	if param.Channel != "" {
		query = query.Where("channel = ?", param.Channel)
	}
	if len(param.TemplateTypes) > 0 {
		query = query.Where("template_type IN (?)", param.TemplateTypes)
	}
	if param.Limit > 0 || param.Offset > 0 {
		query = query.Limit(param.Limit).Offset(param.Offset)
	}
	return query
}

