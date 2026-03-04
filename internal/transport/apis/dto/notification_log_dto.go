package dto

import (
	"strconv"
	"time"

	logEntity "go-boilerplate-clean/internal/entity/notificationlogs"

	"github.com/labstack/echo/v4"
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


// response dto
type NotificationLogResponse struct {
	ID              string     `json:"id"`
	NotificationID  string     `json:"notification_id"`
	UserID          string     `json:"user_id"`
	SendTo          string     `json:"send_to,omitempty"`
	RenderedSubject string     `json:"rendered_subject,omitempty"`
	RenderedMessage string     `json:"rendered_message,omitempty"`
	State           string                 `json:"state,omitempty"`
	RetryCount      int                    `json:"retry_count,omitempty"`
	ErrorMessage    string                 `json:"error_message,omitempty"`
	SentAt          *time.Time             `json:"sent_at,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
}

func (r *NotificationLogResponse) FromEntity(l logEntity.NotificationLog) NotificationLogResponse {
	return NotificationLogResponse{
		ID:              l.ID,
		NotificationID:  l.NotificationID,
		UserID:          l.UserID,
		SendTo:          l.SendTo,
		RenderedSubject: l.RenderedSubject,
		RenderedMessage: l.RenderedMessage,
		// Data:            l.Data,
		State:           l.State.String(),
		RetryCount:      l.RetryCount,
		ErrorMessage:    l.ErrorMessage,
		SentAt:          l.SentAt,
		CreatedAt:       l.CreatedAt,
	}
}
func (r *CreateNotificationLogRequest) ToEntity() logEntity.NotificationLog {
	return logEntity.NotificationLog{
		NotificationID:  r.NotificationID,
		UserID:          r.UserID,
		// Channel:         logEntity.Channel(r.Channel),
		SendTo:          r.SendTo,
		RenderedSubject: r.RenderedSubject,
		RenderedMessage: r.RenderedMessage,
		// Data:            r.Data,
		State:           logEntity.State(r.State),
		RetryCount:      r.RetryCount,
		ErrorMessage:    r.ErrorMessage,
		SentAt:          r.SentAt,
	}
}

// NotificationLogListParamFromQuery mengisi NotificationLogListParam dari query string (untuk GET list).
// Query params: ids, notification_id, user_id, state (comma), limit, offset.
func NotificationLogListParamFromQuery(c echo.Context) *logEntity.NotificationLogListParam {
	param := &logEntity.NotificationLogListParam{}
	if v := c.QueryParam("ids"); v != "" {
		param.IDs = splitTrim(v)
	}
	if v := c.QueryParam("notification_id"); v != "" {
		param.NotificationID = v
	}
	if v := c.QueryParam("user_id"); v != "" {
		param.UserID = v
	}
	if v := c.QueryParam("state"); v != "" {
		param.States = splitTrim(v)
	}
	if v := c.QueryParam("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			param.Limit = n
		}
	}
	if v := c.QueryParam("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			param.Offset = n
		}
	}
	return param
}
