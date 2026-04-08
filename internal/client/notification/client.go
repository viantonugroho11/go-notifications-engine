package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go-boilerplate-clean/internal/entity/notificationinbox"
	notifEntity "go-boilerplate-clean/internal/entity/notifications"
	"go-boilerplate-clean/internal/entity/notificationlogs"
)

// Client memanggil API notification (update, dll.) via HTTP.
type Client interface {
	Update(ctx context.Context, n notifEntity.Notification) (notifEntity.Notification, error)
	CreateInbox(ctx context.Context, n notificationinbox.NotificationInbox) (notificationinbox.NotificationInbox, error)
	UpdateNotificationLog(ctx context.Context, n notificationlogs.NotificationLog) (notificationlogs.NotificationLog, error)
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


func (c *client) CreateInbox(ctx context.Context, n notificationinbox.NotificationInbox) (notificationinbox.NotificationInbox, error) {
	body, err := json.Marshal(notificationInboxToCreateRequest(n))
	if err != nil {
		return notificationinbox.NotificationInbox{}, err
	}
	url := c.baseURL + "/notificationinbox"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return notificationinbox.NotificationInbox{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return notificationinbox.NotificationInbox{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return notificationinbox.NotificationInbox{}, fmt.Errorf("notification API: %s %s", resp.Status, url)
	}
	var out notificationinbox.NotificationInbox
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return notificationinbox.NotificationInbox{}, err
	}
	return out, nil
}

func (c *client) UpdateNotificationLog(ctx context.Context, n notificationlogs.NotificationLog) (notificationlogs.NotificationLog, error) {
body, err := json.Marshal(notificationLogToUpdateRequest(n))
	if err != nil {
		return notificationlogs.NotificationLog{}, err
	}
	url := c.baseURL + "/notificationlogs/" + n.ID
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return notificationlogs.NotificationLog{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return notificationlogs.NotificationLog{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return notificationlogs.NotificationLog{}, fmt.Errorf("notification API: %s %s", resp.Status, url)
	}
	var out notificationlogs.NotificationLog
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return notificationlogs.NotificationLog{}, err
	}
	return out, nil
}