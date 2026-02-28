package model

import "time"

type NotificationLog struct {
	ID              string     `gorm:"column:id;primaryKey"`
	NotificationID  string     `gorm:"column:notification_id;not null"`
	UserID          string     `gorm:"column:user_id;not null"`
	Channel         string     `gorm:"column:channel;not null"`
	SendTo          string     `gorm:"column:send_to;type:text"`
	RenderedSubject string     `gorm:"column:rendered_subject"`
	RenderedMessage string     `gorm:"column:rendered_message;type:text"`
	Data            []byte     `gorm:"column:data;type:jsonb"`
	State           string     `gorm:"column:state;not null"`
	RetryCount      int        `gorm:"column:retry_count"`
	ErrorMessage    string     `gorm:"column:error_message;type:text"`
	SentAt          *time.Time `gorm:"column:sent_at"`
	CreatedAt       time.Time  `gorm:"column:created_at;not null"`
}

func (NotificationLog) TableName() string {
	return "notification_logs"
}
