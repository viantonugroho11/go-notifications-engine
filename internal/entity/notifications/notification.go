package notifications

import (
	"go-boilerplate-clean/internal/entity/notificationlogs"
	"time"
)

type Notification struct {
	ID                     string                             `json:"id"`
	EventKey               string                             `json:"event_key"`
	NotificationTemplateID string                             `json:"notification_template_id"`
	Data                   map[string]interface{}             `json:"data,omitempty"`
	Category               Category                           `json:"category"`
	Channel                Channel                            `json:"channel"`
	State                  string                             `json:"state"`
	ScheduleAt             *time.Time                         `json:"schedule_at,omitempty"`
	CreatedBy              string                             `json:"created_by"`
	UpdatedBy              string                             `json:"updated_by,omitempty"`
	CreatedAt              time.Time                          `json:"created_at"`
	UpdatedAt              *time.Time                         `json:"updated_at,omitempty"`
	NotificationLogs       []notificationlogs.NotificationLog `json:"notification_logs,omitempty"`
}

type Category string

const (
	CategoryPromo         Category = "promo"
	CategoryTransactional Category = "transactional"
	CategorySystem        Category = "system"
	CategoryOther         Category = "other"
)

func (c Category) String() string {
	return string(c)
}

type Channel string

const (
	ChannelEmail    Channel = "email"
	ChannelSMS      Channel = "sms"
	ChannelPush     Channel = "push"
	ChannelWhatsApp Channel = "whatsapp"
	ChannelTelegram Channel = "telegram"
	ChannelLine     Channel = "line"
	ChannelWeChat   Channel = "wechat"
	ChannelWeibo    Channel = "weibo"
	ChannelKakao    Channel = "kakao"
)

type State string

const (
	StateCreated    State = "created"
	StateScheduled  State = "scheduled"
	StateProcessing State = "processing"
	StateRendered   State = "rendered"
	StateSent       State = "sent"
	StateFailed     State = "failed"
	StateCompleted  State = "completed"
)

func (s State) String() string {
	return string(s)
}
