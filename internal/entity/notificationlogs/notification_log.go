package notificationlogs

import "time"

type NotificationLog struct {
	ID               string                 `json:"id"`
	NotificationID   string                 `json:"notification_id"`
	UserID           string                 `json:"user_id"`
	Channel          string                 `json:"channel"`
	SendTo           string                 `json:"send_to,omitempty"`
	RenderedSubject  string                 `json:"rendered_subject,omitempty"`
	RenderedMessage  string                 `json:"rendered_message,omitempty"`
	Data             map[string]interface{} `json:"data,omitempty"`
	State            string                 `json:"state"`
	RetryCount       int                    `json:"retry_count"`
	ErrorMessage     string                 `json:"error_message,omitempty"`
	SentAt           *time.Time             `json:"sent_at,omitempty"`
	CreatedAt        time.Time              `json:"created_at"`
}

const (
	StateQueued     = "queued"
	StateProcessing = "processing"
	StateSent       = "sent"
	StateFailed     = "failed"
	StateCompleted  = "completed"
)
