package model

import "time"

type Notification struct {
	ID                     string    `gorm:"column:id;primaryKey"`
	EventKey               string    `gorm:"column:event_key;not null"`
	NotificationTemplateID string    `gorm:"column:notification_template_id;not null"`
	Data                   []byte    `gorm:"column:data;type:jsonb"`
	Category               string    `gorm:"column:category;not null"`
	State                  string    `gorm:"column:state;not null"`
	ScheduleAt             *time.Time `gorm:"column:schedule_at"`
	CreatedBy              string    `gorm:"column:created_by"`
	UpdatedBy              string    `gorm:"column:updated_by"`
	CreatedAt              time.Time `gorm:"column:created_at;not null"`
	UpdatedAt              *time.Time `gorm:"column:updated_at"`
}

func (Notification) TableName() string {
	return "notifications"
}
