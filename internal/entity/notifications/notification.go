package notifications

import "time"

type Notification struct {
	ID                     string     `json:"id"`
	EventKey               string     `json:"event_key"`
	NotificationTemplateID string     `json:"notification_template_id"`
	Data                   map[string]interface{} `json:"data,omitempty"`
	Category               string     `json:"category"`
	State                  string     `json:"state"`
	ScheduleAt             *time.Time `json:"schedule_at,omitempty"`
	CreatedBy              string     `json:"created_by"`
	UpdatedBy              string     `json:"updated_by,omitempty"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              *time.Time `json:"updated_at,omitempty"`
}

const (
	CategoryPromo         = "promo"
	CategoryTransactional = "transactional"
	CategorySystem        = "system"
	CategoryOther         = "other"
)
