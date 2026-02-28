package notificationinbox

import "time"

type NotificationInbox struct {
	ID                 string     `json:"id"`
	UserID             string     `json:"user_id"`
	NotificationLogID  string     `json:"notification_log_id"`
	IsRead             bool       `json:"is_read"`
	ReadAt             *time.Time `json:"read_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
}
