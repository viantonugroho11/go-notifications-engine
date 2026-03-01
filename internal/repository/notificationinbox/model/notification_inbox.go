package model

import (
	"time"

	notifInboxEntity "go-boilerplate-clean/internal/entity/notificationinbox"
)

type NotificationInbox struct {
	ID                string     `gorm:"column:id;primaryKey"`
	UserID            string     `gorm:"column:user_id;not null"`
	NotificationLogID string     `gorm:"column:notification_log_id;not null"`
	IsRead            bool       `gorm:"column:is_read"`
	ReadAt            *time.Time `gorm:"column:read_at"`
	CreatedAt        time.Time  `gorm:"column:created_at;not null"`
}

func (NotificationInbox) TableName() string {
	return "notification_inbox"
}


func (n NotificationInbox) ToEntity() notifInboxEntity.NotificationInbox {
	return notifInboxEntity.NotificationInbox{
		ID: n.ID,
		UserID: n.UserID,
		NotificationLogID: n.NotificationLogID,
		IsRead: n.IsRead,
		ReadAt: n.ReadAt,
		CreatedAt: n.CreatedAt,
	}
}

func ToDBNotificationInbox(i notifInboxEntity.NotificationInbox) NotificationInbox {
	return NotificationInbox{
		ID: i.ID,
		UserID: i.UserID,
		NotificationLogID: i.NotificationLogID,
		IsRead: i.IsRead,
		ReadAt: i.ReadAt,
		CreatedAt: i.CreatedAt,
	}
}