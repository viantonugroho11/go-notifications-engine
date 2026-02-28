package dto

import (
	"time"

	inboxEntity "go-boilerplate-clean/internal/entity/notificationinbox"
)

type CreateNotificationInboxRequest struct {
	UserID            string     `json:"user_id"`
	NotificationLogID string     `json:"notification_log_id"`
	IsRead            bool       `json:"is_read,omitempty"`
	ReadAt            *time.Time `json:"read_at,omitempty"`
}

type UpdateNotificationInboxRequest struct {
	UserID            string     `json:"user_id"`
	NotificationLogID string     `json:"notification_log_id"`
	IsRead            bool       `json:"is_read"`
	ReadAt            *time.Time `json:"read_at,omitempty"`
}

func (r *CreateNotificationInboxRequest) ToEntity() inboxEntity.NotificationInbox {
	return inboxEntity.NotificationInbox{
		UserID:            r.UserID,
		NotificationLogID: r.NotificationLogID,
		IsRead:            r.IsRead,
		ReadAt:            r.ReadAt,
	}
}
