package model

import "time"

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
