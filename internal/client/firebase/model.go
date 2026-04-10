package firebase

type Notification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Image string `json:"image,omitempty"`
}
type AndroidConfig struct {
	Priority              string `json:"priority,omitempty"`
	CollapseKey           string `json:"collapse_key,omitempty"`
	TimeToLive            int64  `json:"time_to_live,omitempty"`
	RestrictedPackageName string `json:"restricted_package_name,omitempty"`
	PriorityBoost         int64  `json:"priority_boost,omitempty"`
	NotificationChannelId string `json:"notification_channel_id,omitempty"`
}

type WebpushConfig struct {
	Headers      map[string]string `json:"headers,omitempty"`
	Data         map[string]string `json:"data,omitempty"`
	Notification *Notification     `json:"notification,omitempty"`
}

type APNSConfig struct {
	Headers map[string]string `json:"headers,omitempty"`
	Payload map[string]string `json:"payload,omitempty"`
}

type FCMOptions struct {
	AnalyticsLabel string `json:"analytics_label,omitempty"`
}
type Message struct {
	Data         map[string]string `json:"data,omitempty"`
	Notification *Notification     `json:"notification,omitempty"`
	Android      *AndroidConfig    `json:"android,omitempty"`
	Webpush      *WebpushConfig    `json:"webpush,omitempty"`
	APNS         *APNSConfig       `json:"apns,omitempty"`
	FCMOptions   *FCMOptions       `json:"fcm_options,omitempty"`
	Token        string            `json:"token,omitempty"`
	Topic        string            `json:"-"`
	Condition    string            `json:"condition,omitempty"`
}
