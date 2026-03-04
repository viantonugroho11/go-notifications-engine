package dto

import (
	"strconv"
	"time"

	inboxEntity "go-boilerplate-clean/internal/entity/notificationinbox"

	"github.com/labstack/echo/v4"
)

type CreateNotificationInboxRequest struct {
	UserID            string     `json:"user_id"`
	NotificationLogID string     `json:"notification_log_id"`
	IsRead            bool       `json:"is_read,omitempty"`
	ReadAt            *time.Time `json:"read_at,omitempty"`
	Message           string     `json:"message"`
	Subject           string     `json:"subject"`
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
		Message:           r.Message,
		Subject:           r.Subject,
	}
}

// response dto
type NotificationInboxResponse struct {
	ID                string     `json:"id"`
	UserID            string     `json:"user_id"`
	NotificationLogID string     `json:"notification_log_id"`
	IsRead            bool       `json:"is_read"`
	ReadAt            *time.Time `json:"read_at,omitempty"`
	Message           string     `json:"message"`
	Subject           string     `json:"subject"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
}

func (r *NotificationInboxResponse) FromEntity(i inboxEntity.NotificationInbox) NotificationInboxResponse {
	return NotificationInboxResponse{
		ID:                i.ID,
		UserID:            i.UserID,
		NotificationLogID: i.NotificationLogID,
		Message:           i.Message,
		Subject:           i.Subject,
		IsRead:            i.IsRead,
		ReadAt:            i.ReadAt,
		CreatedAt:         i.CreatedAt,
	}
}

// NotificationInboxListParamFromQuery mengisi NotificationInboxListParam dari query string (untuk GET list).
// Query params: ids, user_id, notification_log_id, is_read, limit, offset.
func NotificationInboxListParamFromQuery(c echo.Context) *inboxEntity.NotificationInboxListParam {
	param := &inboxEntity.NotificationInboxListParam{}
	if v := c.QueryParam("ids"); v != "" {
		param.IDs = splitTrim(v)
	}
	if v := c.QueryParam("user_id"); v != "" {
		param.UserID = v
	}
	if v := c.QueryParam("notification_log_id"); v != "" {
		param.NotificationLogID = v
	}
	if v := c.QueryParam("is_read"); v != "" {
		b, err := strconv.ParseBool(v)
		if err == nil {
			param.IsRead = &b
		}
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
