package notifications

import (
	"bytes"
	"go-boilerplate-clean/internal/entity/notificationlogs"
	"text/template"
	"time"
)

type Notification struct {
	ID                     string                             `json:"id"`
	EventKey               string                             `json:"event_key"`
	NotificationTemplateID string                             `json:"notification_template_id"`
	Data                   map[string]interface{}             `json:"data,omitempty"`
	Category               Category                           `json:"category"`
	Channel                Channel                            `json:"channel"`
	State                  State                              `json:"state"`
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

func (c Channel) String() string {
	return string(c)
}

type State string

// -- created / scheduled / processing / sent / failed / completed
const (
	StateCreated    State = "CREATED"
	StateScheduled  State = "SCHEDULED"
	StateProcessing State = "PROCESSING"
	StateSent       State = "SENT"
	StateFailed     State = "FAILED"
	StateCompleted  State = "COMPLETED"
)

func (s State) String() string {
	return string(s)
}

// NotificationListParam dipakai untuk filter dan pagination List notification.
type NotificationListParam struct {
	IDs                    []string
	EventKey               string
	NotificationTemplateID string
	Channel                string
	Categories             []string
	States                 []string
	Page, Limit, Offset    int
}

// message template
func (n *Notification) GenerateRenderedMessage(messageTemplate string) string {
	t, err := template.New("tmpl").
		Option("missingkey=zero").
		Parse(messageTemplate)

	if err != nil {
		return ""
	}

	var buf bytes.Buffer

	err = t.Execute(&buf, n.Data)
	if err != nil {
		return ""
	}

	return buf.String()
}

type NotificationProducerMessage struct {
	NotificationID         string                             `json:"notification_id"`
	EventKey               string                             `json:"event_key"`
	NotificationTemplateID string                             `json:"notification_template_id"`
	Data                   map[string]interface{}             `json:"data,omitempty"`
	Channel                string                             `json:"channel"`
	Category               string                             `json:"category"`
	State                  string                             `json:"state"`
	ScheduleAt             *time.Time                         `json:"schedule_at,omitempty"`
	CreatedBy              string                             `json:"created_by,omitempty"`
	UpdatedBy              string                             `json:"updated_by,omitempty"`
	NotificationLogs       []notificationlogs.NotificationLog `json:"notification_logs,omitempty"`
}

func (n *Notification) ToProducerMessage() NotificationProducerMessage {
	return NotificationProducerMessage{
		NotificationID:         n.ID,
		EventKey:               n.EventKey,
		NotificationTemplateID: n.NotificationTemplateID,
		Data:                   n.Data,
		Channel:                n.Channel.String(),
		Category:               n.Category.String(),
		State:                  n.State.String(),
		ScheduleAt:             n.ScheduleAt,
		CreatedBy:              n.CreatedBy,
		UpdatedBy:              n.UpdatedBy,
		NotificationLogs:       n.NotificationLogs,
	}
}

// ToNotification mengonversi NotificationProducerMessage (dari event Kafka) ke entity Notification untuk update.
func (m *NotificationProducerMessage) ToNotification() Notification {
	return Notification{
		ID:                     m.NotificationID,
		EventKey:               m.EventKey,
		NotificationTemplateID: m.NotificationTemplateID,
		Data:                   m.Data,
		Channel:                Channel(m.Channel),
		Category:               Category(m.Category),
		State:                  State(m.State),
		ScheduleAt:             m.ScheduleAt,
		CreatedBy:              m.CreatedBy,
		UpdatedBy:              m.UpdatedBy,
		NotificationLogs:       m.NotificationLogs,
	}
}

type NotificationsEventMessage struct {
	Action string                      `json:"action"`
	After  NotificationProducerMessage `json:"after"`
	Before NotificationProducerMessage `json:"before"`
}

func ToNotificationsEventMessage(action string, after NotificationProducerMessage, before NotificationProducerMessage) NotificationsEventMessage {
	return NotificationsEventMessage{
		Action: action,
		After:  after,
		Before: before,
	}
}

type NotificationEventUsecase struct {
	ID                     string                           `json:"id"`
	EventKey               string                           `json:"event_key"`
	NotificationTemplateID string                           `json:"notification_template_id"`
	Data                   map[string]interface{}           `json:"data,omitempty"`
	Category               Category                         `json:"category"`
	Channel                Channel                          `json:"channel"`
	State                  State                            `json:"state"`
	ScheduleAt             *time.Time                       `json:"schedule_at,omitempty"`
	CreatedBy              string                           `json:"created_by"`
	UpdatedBy              string                           `json:"updated_by,omitempty"`
	CreatedAt              time.Time                        `json:"created_at"`
	UpdatedAt              *time.Time                       `json:"updated_at,omitempty"`
	NotificationLogs       notificationlogs.NotificationLog `json:"notification_logs,omitempty"`
}
