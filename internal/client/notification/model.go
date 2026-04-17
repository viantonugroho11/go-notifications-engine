package notification

import (
	"time"

	"github.com/viantonugroho11/go-notifications-engine/internal/entity/notificationinbox"
	"github.com/viantonugroho11/go-notifications-engine/internal/entity/notificationlogs"
	notifEntity "github.com/viantonugroho11/go-notifications-engine/internal/entity/notifications"
)

type updateRequest struct {
	EventKey               string                 `json:"event_key"`
	NotificationTemplateID string                 `json:"notification_template_id"`
	Data                   map[string]interface{} `json:"data,omitempty"`
	Channel                string                 `json:"channel"`
	Category               string                 `json:"category"`
	State                  string                 `json:"state"`
	ScheduleAt             *string                `json:"schedule_at,omitempty"`
	UpdatedBy              string                 `json:"updated_by,omitempty"`
}

func notificationToUpdateRequest(n notifEntity.Notification) updateRequest {
	r := updateRequest{
		EventKey:               n.EventKey,
		NotificationTemplateID: n.NotificationTemplateID,
		Data:                   n.Data,
		Channel:                n.Channel.String(),
		Category:               n.Category.String(),
		State:                  n.State.String(),
		UpdatedBy:              n.UpdatedBy,
	}
	if n.ScheduleAt != nil {
		s := n.ScheduleAt.Format("2006-01-02T15:04:05Z07:00")
		r.ScheduleAt = &s
	}
	return r
}

type createInboxRequest struct {
	UserID            string     `json:"user_id"`
	NotificationLogID string     `json:"notification_log_id"`
	IsRead            bool       `json:"is_read,omitempty"`
	ReadAt            *time.Time `json:"read_at,omitempty"`
	Message           string     `json:"message"`
	Subject           string     `json:"subject"`
}

func notificationInboxToCreateRequest(n notificationinbox.NotificationInbox) createInboxRequest {
	r := createInboxRequest{
		UserID:            n.UserID,
		NotificationLogID: n.NotificationLogID,
		IsRead:            false,
		Message:           n.Message,
		Subject:           n.Subject,
	}
	if n.ReadAt != nil {
		r.ReadAt = n.ReadAt
		r.IsRead = true
	}
	return r
}

type updateNotificationLogRequest struct {
	State string `json:"state"`
}

func notificationLogToUpdateRequest(n notificationlogs.NotificationLog) updateNotificationLogRequest {
	return updateNotificationLogRequest{
		State: n.State.String(),
	}
}
