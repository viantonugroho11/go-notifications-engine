package dto

import (
	"strconv"
	"strings"
	"time"

	notificationlogs "go-boilerplate-clean/internal/entity/notificationlogs"
	"go-boilerplate-clean/internal/entity/notifications"
	notifEntity "go-boilerplate-clean/internal/entity/notifications"

	"github.com/labstack/echo/v4"
)

type CreateNotificationRequest struct {
	EventKey               string         `json:"event_key"`
	NotificationTemplateID string         `json:"notification_template_id"`
	Data                   map[string]any `json:"data,omitempty"`
	Channel                string         `json:"channel"`
	Category               string         `json:"category"`
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
		State:                  notifEntity.StateCreated.String(),
		ScheduleAt:             r.ScheduleAt,
		CreatedBy:              r.CreatedBy,
		Channel:                notifications.Channel(r.Channel),
		NotificationLogs:       logs,
	}
}

type 	NotificationResponse struct {
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

func FromUsecaseNotification(n notifEntity.Notification) NotificationResponse {
	logs := []notificationlogs.NotificationLog{}
	for _, log := range n.NotificationLogs {
		logs = append(logs, notificationlogs.NotificationLog{
			ID:        log.ID,
			UserID:    log.UserID,
			State:     log.State,
			CreatedAt: log.CreatedAt,
		})
	}
	return NotificationResponse{
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

// NotificationListParamFromQuery mengisi NotificationListParam dari query string (untuk GET list).
// Query params: event_key, notification_template_id, channel, category (comma), state (comma), ids (comma), page, limit, offset.
func NotificationListParamFromQuery(c echo.Context) *notifEntity.NotificationListParam {
	param := &notifEntity.NotificationListParam{}
	if v := c.QueryParam("event_key"); v != "" {
		param.EventKey = v
	}
	if v := c.QueryParam("notification_template_id"); v != "" {
		param.NotificationTemplateID = v
	}
	if v := c.QueryParam("channel"); v != "" {
		param.Channel = v
	}
	if v := c.QueryParam("category"); v != "" {
		param.Categories = splitTrim(v)
	}
	if v := c.QueryParam("state"); v != "" {
		param.States = splitTrim(v)
	}
	if v := c.QueryParam("ids"); v != "" {
		param.IDs = splitTrim(v)
	}
	if v := c.QueryParam("page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			param.Page = n
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

func splitTrim(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}

// parameter dto
type NotificationParam struct {
	ID                     string `param:"id"`
	EventKey               string `param:"event_key"`
	NotificationTemplateID string `param:"notification_template_id"`
	Channel                string `param:"channel"`
	Category               string `param:"category"`
	State                  string `param:"state"`
	Limit                  int    `param:"limit"`
	Offset                 int    `param:"offset"`
	Sort                   string `param:"sort"`
	SortBy                 string `param:"sort_by"`
}

type NotificationListResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	Total         int                    `json:"total"`
	Page          int                    `json:"page"`
	Limit         int                    `json:"limit"`
	Offset        int                    `json:"offset"`
}



