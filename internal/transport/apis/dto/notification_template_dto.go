package dto

import (
	"strconv"
	"time"

	tplEntity "github.com/viantonugroho11/go-notifications-engine/internal/entity/notificationtemplates"

	"github.com/labstack/echo/v4"
)

type CreateNotificationTemplateRequest struct {
	Name          string                 `json:"name"`
	Subject       string                 `json:"subject,omitempty"`
	Body          string                 `json:"body,omitempty"`
	PayloadSchema map[string]interface{} `json:"payload_schema,omitempty"`
	Channel       string                 `json:"channel"`
	TemplateType  string                 `json:"template_type,omitempty"`
}

// response dto
type NotificationTemplateResponse struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	Subject       string         `json:"subject,omitempty"`
	Body          string         `json:"body,omitempty"`
	PayloadSchema map[string]any `json:"payload_schema,omitempty"`
	Channel       string         `json:"channel"`
	TemplateType  string         `json:"template_type,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     *time.Time     `json:"updated_at,omitempty"`
	DeletedAt     *time.Time     `json:"deleted_at,omitempty"`
}

func (r *NotificationTemplateResponse) FromEntity(t tplEntity.NotificationTemplate) NotificationTemplateResponse {
	return NotificationTemplateResponse{
		ID:            t.ID,
		Name:          t.Name,
		Subject:       t.Subject,
		Body:          t.Body,
		PayloadSchema: t.PayloadSchema,
		Channel:       t.Channel,
		TemplateType:  t.TemplateType.String(),
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
		DeletedAt:     t.DeletedAt,
	}
}

type UpdateNotificationTemplateRequest struct {
	Name          string                 `json:"name"`
	Subject       string                 `json:"subject,omitempty"`
	Body          string                 `json:"body,omitempty"`
	PayloadSchema map[string]interface{} `json:"payload_schema,omitempty"`
	Channel       string                 `json:"channel"`
	TemplateType  string                 `json:"template_type,omitempty"`
}

func (r *CreateNotificationTemplateRequest) ToEntity() tplEntity.NotificationTemplate {
	return tplEntity.NotificationTemplate{
		Name:          r.Name,
		Subject:       r.Subject,
		Body:          r.Body,
		PayloadSchema: r.PayloadSchema,
		Channel:       r.Channel,
		TemplateType:  tplEntity.TemplateType(r.TemplateType),
	}
}

// update to entity,
func (r *UpdateNotificationTemplateRequest) ToEntity() tplEntity.NotificationTemplate {
	return tplEntity.NotificationTemplate{
		Name:          r.Name,
		Subject:       r.Subject,
		Body:          r.Body,
		PayloadSchema: r.PayloadSchema,
		Channel:       r.Channel,
		TemplateType:  tplEntity.TemplateType(r.TemplateType),
	}
}

// NotificationTemplateListParamFromQuery mengisi NotificationTemplateListParam dari query string (untuk GET list).
// Query params: ids, name, channel, template_type (comma), limit, offset.
func NotificationTemplateListParamFromQuery(c echo.Context) *tplEntity.NotificationTemplateListParam {
	param := &tplEntity.NotificationTemplateListParam{}
	if v := c.QueryParam("ids"); v != "" {
		param.IDs = splitTrim(v)
	}
	if v := c.QueryParam("name"); v != "" {
		param.Name = v
	}
	if v := c.QueryParam("channel"); v != "" {
		param.Channel = v
	}
	if v := c.QueryParam("template_type"); v != "" {
		param.TemplateTypes = splitTrim(v)
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
