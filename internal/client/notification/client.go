package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	notifEntity "go-boilerplate-clean/internal/entity/notifications"
)

// Client memanggil API notification (update, dll.) via HTTP.
type Client interface {
	Update(ctx context.Context, n notifEntity.Notification) (notifEntity.Notification, error)
}

type client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient membuat notification API client.
func NewClient(baseURL string, httpClient *http.Client) Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &client{baseURL: baseURL, httpClient: httpClient}
}

// Update memanggil PUT /notifications/:id.
func (c *client) Update(ctx context.Context, n notifEntity.Notification) (notifEntity.Notification, error) {
	body, err := json.Marshal(notificationToUpdateRequest(n))
	if err != nil {
		return notifEntity.Notification{}, err
	}
	url := c.baseURL + "/notifications/" + n.ID
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return notifEntity.Notification{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return notifEntity.Notification{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return notifEntity.Notification{}, fmt.Errorf("notification API: %s %s", resp.Status, url)
	}
	var out notifEntity.Notification
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return notifEntity.Notification{}, err
	}
	return out, nil
}

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
