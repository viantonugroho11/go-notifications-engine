package notificationinbox

import "time"

type NotificationInbox struct {
	ID                 string     `json:"id"`
	UserID             string     `json:"user_id"`
	NotificationLogID  string     `json:"notification_log_id"`
	Message            string     `json:"message"`
	Subject            string     `json:"subject"`
	IsRead             bool       `json:"is_read"`
	ReadAt             *time.Time `json:"read_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
}

// NotificationInboxListParam dipakai untuk filter dan pagination List notification inbox.
type NotificationInboxListParam struct {
	IDs                []string
	UserID             string
	NotificationLogID  string
	IsRead             *bool
	Limit, Offset      int
}
