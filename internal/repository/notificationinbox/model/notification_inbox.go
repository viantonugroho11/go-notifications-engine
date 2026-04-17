package model

import (
	"time"

	notifInboxEntity "github.com/viantonugroho11/go-notifications-engine/internal/entity/notificationinbox"

	"gorm.io/gorm"
)

type NotificationInbox struct {
	ID                string     `gorm:"column:id;primaryKey"`
	UserID            string     `gorm:"column:user_id;not null"`
	NotificationLogID string     `gorm:"column:notification_log_id;not null"`
	IsRead            bool       `gorm:"column:is_read"`
	ReadAt            *time.Time `gorm:"column:read_at"`
	CreatedAt         time.Time  `gorm:"column:created_at;not null"`
	Message           string     `gorm:"column:message;type:text"`
	Subject           string     `gorm:"column:subject;type:text"`
}

func (NotificationInbox) TableName() string {
	return "notification_inbox"
}

func (n NotificationInbox) ToEntity() notifInboxEntity.NotificationInbox {
	return notifInboxEntity.NotificationInbox{
		ID:                n.ID,
		UserID:            n.UserID,
		NotificationLogID: n.NotificationLogID,
		IsRead:            n.IsRead,
		ReadAt:            n.ReadAt,
		CreatedAt:         n.CreatedAt,
		Message:           n.Message,
		Subject:           n.Subject,
	}
}

func ToDBNotificationInbox(i notifInboxEntity.NotificationInbox) NotificationInbox {
	return NotificationInbox{
		ID:                i.ID,
		UserID:            i.UserID,
		NotificationLogID: i.NotificationLogID,
		IsRead:            i.IsRead,
		ReadAt:            i.ReadAt,
		CreatedAt:         i.CreatedAt,
		Message:           i.Message,
		Subject:           i.Subject,
	}
}

// ApplyListParam menerapkan filter dan pagination dari param ke query. Jika param nil, query tidak diubah.
func ApplyListParam(query *gorm.DB, param *notifInboxEntity.NotificationInboxListParam) *gorm.DB {
	if param == nil {
		return query
	}
	if len(param.IDs) > 0 {
		query = query.Where("id IN (?)", param.IDs)
	}
	if param.UserID != "" {
		query = query.Where("user_id = ?", param.UserID)
	}
	if param.NotificationLogID != "" {
		query = query.Where("notification_log_id = ?", param.NotificationLogID)
	}
	if param.IsRead != nil {
		query = query.Where("is_read = ?", *param.IsRead)
	}
	if param.Limit > 0 || param.Offset > 0 {
		query = query.Limit(param.Limit).Offset(param.Offset)
	}
	return query
}
