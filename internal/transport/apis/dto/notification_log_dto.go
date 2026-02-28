package dto

import (
	"time"

	logEntity "go-boilerplate-clean/internal/entity/notificationlogs"
)

type CreateNotificationLogRequest struct {
	NotificationID   string                 `json:"notification_id"`
	UserID           string                 `json:"user_id"`
	Channel          string                 `json:"channel"`
	SendTo           string                 `json:"send_to,omitempty"`
	RenderedSubject  string                 `json:"rendered_subject,omitempty"`
	RenderedMessage  string                 `json:"rendered_message,omitempty"`
	Data             map[string]interface{} `json:"data,omitempty"`
	State            string                 `json:"state,omitempty"`
	RetryCount       int                    `json:"retry_count,omitempty"`
	ErrorMessage     string                 `json:"error_message,omitempty"`
	SentAt           *time.Time             `json:"sent_at,omitempty"`
}

type UpdateNotificationLogRequest struct {
	NotificationID   string                 `json:"notification_id"`
	UserID           string                 `json:"user_id"`
	Channel          string                 `json:"channel"`
	SendTo           string                 `json:"send_to,omitempty"`
	RenderedSubject  string                 `json:"rendered_subject,omitempty"`
	RenderedMessage  string                 `json:"rendered_message,omitempty"`
	Data             map[string]interface{} `json:"data,omitempty"`
	State            string                 `json:"state"`
	RetryCount       int                    `json:"retry_count,omitempty"`
	ErrorMessage     string                 `json:"error_message,omitempty"`
	SentAt           *time.Time             `json:"sent_at,omitempty"`
}

func (r *CreateNotificationLogRequest) ToEntity() logEntity.NotificationLog {
	return logEntity.NotificationLog{
		NotificationID:  r.NotificationID,
		UserID:          r.UserID,
		Channel:         r.Channel,
		SendTo:          r.SendTo,
		RenderedSubject: r.RenderedSubject,
		RenderedMessage: r.RenderedMessage,
		Data:            r.Data,
		State:           r.State,
		RetryCount:      r.RetryCount,
		ErrorMessage:    r.ErrorMessage,
		SentAt:          r.SentAt,
	}
}
