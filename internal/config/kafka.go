package config

type Kafka struct {
	Brokers   []string `json:"brokers"`
	ClientID  string   `json:"client_id"`
	GroupID   string   `json:"group_id"`
	Topic     string   `json:"topic"`
	TopicSent string   `json:"topic_sent"` // topic untuk event "sent" (kirim via email/firebase)
}