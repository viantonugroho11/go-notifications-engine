package model

import (
	"encoding/json"
	"go-boilerplate-clean/internal/entity/notificationlogs"
	notifEntity "go-boilerplate-clean/internal/entity/notifications"
	notifLogModel "go-boilerplate-clean/internal/repository/notificationlog/model"
	"time"
)

type Notification struct {
	ID                     string                          `gorm:"column:id;primaryKey"`
	EventKey               string                          `gorm:"column:event_key;not null"`
	NotificationTemplateID string                          `gorm:"column:notification_template_id;not null"`
	Data                   []byte                          `gorm:"column:data;type:jsonb"`
	Category               string                          `gorm:"column:category;not null"`
	State                  string                          `gorm:"column:state;not null"`
	Channel                string                          `gorm:"column:channel;not null"`
	ScheduleAt             *time.Time                      `gorm:"column:schedule_at"`
	CreatedBy              string                          `gorm:"column:created_by"`
	UpdatedBy              string                          `gorm:"column:updated_by"`
	CreatedAt              time.Time                       `gorm:"column:created_at;not null"`
	UpdatedAt              *time.Time                      `gorm:"column:updated_at"`
	NotificationLogs       []notifLogModel.NotificationLog `gorm:"foreignKey:notification_id"`
}

func (Notification) TableName() string {
	return "notifications"
}

func (n Notification) ToEntity() notifEntity.Notification {
	var data map[string]any
	if len(n.Data) > 0 {
		_ = json.Unmarshal(n.Data, &data)
	}
	return notifEntity.Notification{
		ID:                     n.ID,
		EventKey:               n.EventKey,
		NotificationTemplateID: n.NotificationTemplateID,
		Data:                   data,
		Category:               notifEntity.Category(n.Category),
		State:                  n.State,
		Channel:                notifEntity.Channel(n.Channel),
		ScheduleAt:             n.ScheduleAt,
		CreatedBy:              n.CreatedBy,
		UpdatedBy:              n.UpdatedBy,
		CreatedAt:              n.CreatedAt,
		UpdatedAt:              n.UpdatedAt,
		NotificationLogs:       ToEntityNotificationLogs(n.NotificationLogs),
	}
}

func ToDBNotification(n notifEntity.Notification) Notification {
	dataJSON, _ := json.Marshal(n.Data)
	return Notification{
		ID:                     n.ID,
		EventKey:               n.EventKey,
		NotificationTemplateID: n.NotificationTemplateID,
		Data:                   dataJSON,
		Category:               n.Category.String(),
		State:                  n.State,
		ScheduleAt:             n.ScheduleAt,
		CreatedBy:              n.CreatedBy,
		UpdatedBy:              n.UpdatedBy,
		CreatedAt:              n.CreatedAt,
		UpdatedAt:              n.UpdatedAt,
		NotificationLogs:       ToDBNotificationLogs(n.NotificationLogs),
	}
}

func ToDBNotificationLogs(logs []notificationlogs.NotificationLog) []notifLogModel.NotificationLog {
	var notificationLogs []notifLogModel.NotificationLog
	for _, l := range logs {
		notificationLogs = append(notificationLogs, notifLogModel.ToDBNotificationLog(l))
	}
	return notificationLogs
}

func ToEntityNotificationLogs(logs []notifLogModel.NotificationLog) []notificationlogs.NotificationLog {
	var notificationLogs []notificationlogs.NotificationLog
	for _, l := range logs {
		notificationLogs = append(notificationLogs, l.ToEntity())
	}
	return notificationLogs
}
