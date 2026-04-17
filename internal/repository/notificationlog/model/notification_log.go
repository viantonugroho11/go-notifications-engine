package model

import (
	// "encoding/json"
	"time"

	"github.com/viantonugroho11/go-notifications-engine/internal/entity/notificationlogs"

	"gorm.io/gorm"
)

type NotificationLog struct {
	ID             string `gorm:"column:id;primaryKey"`
	NotificationID string `gorm:"column:notification_id;not null"`
	UserID         string `gorm:"column:user_id;not null"`
	// Channel         string     `gorm:"column:channel;not null"`
	SendTo          string `gorm:"column:send_to;type:text"`
	RenderedSubject string `gorm:"column:rendered_subject"`
	RenderedMessage string `gorm:"column:rendered_message;type:text"`
	// Data            []byte     `gorm:"column:data;type:jsonb"`
	State        string     `gorm:"column:state;not null"`
	RetryCount   int        `gorm:"column:retry_count"`
	ErrorMessage string     `gorm:"column:error_message;type:text"`
	SentAt       *time.Time `gorm:"column:sent_at"`
	CreatedAt    time.Time  `gorm:"column:created_at;not null"`
}

func (NotificationLog) TableName() string {
	return "notification_logs"
}

func (n NotificationLog) ToEntity() notificationlogs.NotificationLog {
	return notificationlogs.NotificationLog{
		ID:             n.ID,
		NotificationID: n.NotificationID,
		UserID:         n.UserID,
		// Channel: notificationlogs.Channel(n.Channel),
		SendTo:          n.SendTo,
		RenderedSubject: n.RenderedSubject,
		RenderedMessage: n.RenderedMessage,
		// Data: data,
		State:        notificationlogs.State(n.State),
		RetryCount:   n.RetryCount,
		ErrorMessage: n.ErrorMessage,
		SentAt:       n.SentAt,
		CreatedAt:    n.CreatedAt,
	}
}

func ToDBNotificationLog(l notificationlogs.NotificationLog) NotificationLog {
	// dataJSON, _ := json.Marshal(l.Data)
	return NotificationLog{
		ID:             l.ID,
		NotificationID: l.NotificationID,
		UserID:         l.UserID,
		// Channel: l.Channel.String(),
		SendTo:          l.SendTo,
		RenderedSubject: l.RenderedSubject,
		RenderedMessage: l.RenderedMessage,
		// Data: dataJSON,
		State:        l.State.String(),
		RetryCount:   l.RetryCount,
		ErrorMessage: l.ErrorMessage,
		SentAt:       l.SentAt,
		CreatedAt:    l.CreatedAt,
	}
}

// ApplyListParam menerapkan filter dan pagination dari param ke query. Jika param nil, query tidak diubah.
func ApplyListParam(query *gorm.DB, param *notificationlogs.NotificationLogListParam) *gorm.DB {
	if param == nil {
		return query
	}
	if len(param.IDs) > 0 {
		query = query.Where("id IN (?)", param.IDs)
	}
	if param.NotificationID != "" {
		query = query.Where("notification_id = ?", param.NotificationID)
	}
	if param.UserID != "" {
		query = query.Where("user_id = ?", param.UserID)
	}
	if len(param.States) > 0 {
		query = query.Where("state IN (?)", param.States)
	}
	if param.Limit > 0 || param.Offset > 0 {
		query = query.Limit(param.Limit).Offset(param.Offset)
	}
	return query
}
