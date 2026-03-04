package notificationlogs

import "time"

type NotificationLog struct {
	ID               string                 `json:"id"`
	NotificationID   string                 `json:"notification_id"`
	UserID           string                 `json:"user_id"`
	// Channel          Channel                 `json:"channel"`
	SendTo           string                 `json:"send_to,omitempty"`
	RenderedSubject  string                 `json:"rendered_subject,omitempty"`
	RenderedMessage  string                 `json:"rendered_message,omitempty"`
	// Data             map[string]interface{} `json:"data,omitempty"`
	State            State                  `json:"state"`
	RetryCount       int                    `json:"retry_count"`
	ErrorMessage     string                 `json:"error_message,omitempty"`
	SentAt           *time.Time             `json:"sent_at,omitempty"`
	CreatedAt        time.Time              `json:"created_at"`
}


type State string

const (
	StateQueued     State = "queued"
	StateProcessing State = "processing"
	StateSent       State = "sent"
	StateFailed     State = "failed"
	StateCompleted  State = "completed"
)

func (s State) String() string {
	return string(s)
}

type Channel string

const (
	ChannelEmail Channel = "email"
	ChannelSMS   Channel = "sms"
	ChannelPush  Channel = "push"
	ChannelWhatsApp Channel = "whatsapp"
	ChannelTelegram Channel = "telegram"
	ChannelLine Channel = "line"
	ChannelWeChat Channel = "wechat"
	ChannelWeibo Channel = "weibo"
	ChannelKakao Channel = "kakao"
)

func (c Channel) String() string {
	return string(c)
}

// NotificationLogListParam dipakai untuk filter dan pagination List notification log.
type NotificationLogListParam struct {
	IDs            []string
	NotificationID string
	UserID         string
	States         []string
	Limit, Offset  int
}