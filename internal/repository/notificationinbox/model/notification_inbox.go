package model

import "time"

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
