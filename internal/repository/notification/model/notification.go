package model

import (
	"encoding/json"
	"time"

	"github.com/viantonugroho11/go-notifications-engine/internal/entity/notificationlogs"
	notifEntity "github.com/viantonugroho11/go-notifications-engine/internal/entity/notifications"
	notifLogModel "github.com/viantonugroho11/go-notifications-engine/internal/repository/notificationlog/model"

	"gorm.io/gorm"
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
		State:                  notifEntity.State(n.State),
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
		State:                  notifEntity.State(n.State).String(),
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

// ApplyListParam menerapkan filter dan pagination dari param ke query. Jika param nil, query tidak diubah.
func ApplyListParam(query *gorm.DB, param *notifEntity.NotificationListParam) *gorm.DB {
	if param == nil {
		return query
	}
	if len(param.IDs) > 0 {
		query = query.Where("id IN (?)", param.IDs)
	}
	if param.EventKey != "" {
		query = query.Where("event_key = ?", param.EventKey)
	}
	if param.NotificationTemplateID != "" {
		query = query.Where("notification_template_id = ?", param.NotificationTemplateID)
	}
	if param.Channel != "" {
		query = query.Where("channel = ?", param.Channel)
	}
	if len(param.Categories) > 0 {
		query = query.Where("category IN (?)", param.Categories)
	}
	if len(param.States) > 0 {
		query = query.Where("state IN (?)", param.States)
	}
	if param.Limit > 0 || param.Offset > 0 {
		query = query.Limit(param.Limit).Offset(param.Offset)
	}
	return query
}
