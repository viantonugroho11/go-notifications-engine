package dto

import (
	"time"

	notifEntity "go-boilerplate-clean/internal/entity/notifications"
)

type CreateNotificationRequest struct {
	EventKey               string                 `json:"event_key"`
	NotificationTemplateID string                 `json:"notification_template_id"`
	Data                   map[string]interface{} `json:"data,omitempty"`
	Category               string                 `json:"category"`
	State                  string                 `json:"state"`
	ScheduleAt             *time.Time             `json:"schedule_at,omitempty"`
	CreatedBy               string                 `json:"created_by,omitempty"`
}

type UpdateNotificationRequest struct {
	EventKey               string                 `json:"event_key"`
	NotificationTemplateID string                 `json:"notification_template_id"`
	Data                   map[string]interface{} `json:"data,omitempty"`
	Category               string                 `json:"category"`
	State                  string                 `json:"state"`
	ScheduleAt             *time.Time             `json:"schedule_at,omitempty"`
	UpdatedBy              string                 `json:"updated_by,omitempty"`
}

func (r *CreateNotificationRequest) ToEntity() notifEntity.Notification {
	return notifEntity.Notification{
		EventKey:               r.EventKey,
		NotificationTemplateID: r.NotificationTemplateID,
		Data:                   r.Data,
		Category:               r.Category,
		State:                  r.State,
		ScheduleAt:             r.ScheduleAt,
		CreatedBy:              r.CreatedBy,
	}
}
