package dto

import (
	"time"

	notificationlogs "go-boilerplate-clean/internal/entity/notificationlogs"
	"go-boilerplate-clean/internal/entity/notifications"
	notifEntity "go-boilerplate-clean/internal/entity/notifications"
)

type CreateNotificationRequest struct {
	EventKey               string         `json:"event_key"`
	NotificationTemplateID string         `json:"notification_template_id"`
	Data                   map[string]any `json:"data,omitempty"`
	Channel                string         `json:"channel"`
	Category               string         `json:"category"`
	State                  string         `json:"state"`
	ScheduleAt             *time.Time     `json:"schedule_at,omitempty"`
	CreatedBy              string         `json:"created_by,omitempty"`
	UserIDs                []string       `json:"user_ids"`
}

type UpdateNotificationRequest struct {
	EventKey               string         `json:"event_key"`
	NotificationTemplateID string         `json:"notification_template_id"`
	Data                   map[string]any `json:"data,omitempty"`
	Channel                string         `json:"channel"`
	Category               string         `json:"category"`
	State                  string         `json:"state"`
	ScheduleAt             *time.Time     `json:"schedule_at,omitempty"`
	UpdatedBy              string         `json:"updated_by,omitempty"`
}

func (r *CreateNotificationRequest) ToEntity() notifEntity.Notification {
	logs := []notificationlogs.NotificationLog{}
	for _, userID := range r.UserIDs {
		logs = append(logs, notificationlogs.NotificationLog{
			UserID: userID,
			State:  notificationlogs.StateQueued,
		})
	}
	return notifEntity.Notification{
		EventKey:               r.EventKey,
		NotificationTemplateID: r.NotificationTemplateID,
		Data:                   r.Data,
		Category:               notifications.Category(r.Category),
		State:                  r.State,
		ScheduleAt:             r.ScheduleAt,
		CreatedBy:              r.CreatedBy,
		Channel:                notifications.Channel(r.Channel),
		NotificationLogs:       logs,
	}
}

type CreateNotificationResponse struct {
	ID                     string                             `json:"id"`
	EventKey               string                             `json:"event_key"`
	NotificationTemplateID string                             `json:"notification_template_id"`
	Data                   map[string]any                     `json:"data,omitempty"`
	Channel                notifications.Channel              `json:"channel"`
	Category               notifications.Category             `json:"category"`
	State                  string                             `json:"state"`
	ScheduleAt             *time.Time                         `json:"schedule_at,omitempty"`
	CreatedBy              string                             `json:"created_by,omitempty"`
	UpdatedBy              string                             `json:"updated_by,omitempty"`
	CreatedAt              time.Time                          `json:"created_at"`
	UpdatedAt              *time.Time                         `json:"updated_at,omitempty"`
	NotificationLogs       []notificationlogs.NotificationLog `json:"notification_logs"`
}

func FromUsecaseNotification(n notifEntity.Notification) CreateNotificationResponse {
	logs := []notificationlogs.NotificationLog{}
	for _, log := range n.NotificationLogs {
		logs = append(logs, notificationlogs.NotificationLog{
			ID:        log.ID,
			UserID:    log.UserID,
			State:     log.State,
			CreatedAt: log.CreatedAt,
		})
	}
	return CreateNotificationResponse{
		ID:                     n.ID,
		EventKey:               n.EventKey,
		NotificationTemplateID: n.NotificationTemplateID,
		Data:                   n.Data,
		Channel:                n.Channel,
		Category:               notifications.Category(n.Category),
		State:                  n.State,
		ScheduleAt:             n.ScheduleAt,
		CreatedBy:              n.CreatedBy,
		UpdatedBy:              n.UpdatedBy,
		CreatedAt:              n.CreatedAt,
		UpdatedAt:              n.UpdatedAt,
		NotificationLogs:       logs,
	}
}
