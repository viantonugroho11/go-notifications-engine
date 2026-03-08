CREATE TABLE notification_logs (
    id UUID PRIMARY KEY,
    notification_id UUID NOT NULL, -- notification that the log is for
    user_id UUID NOT NULL, -- user who received the notification

    -- channel VARCHAR(20) NOT NULL, // email / sms / push / whatsapp / telegram / line / wechat / weibo / kakao
    send_to TEXT, -- email address or device token or phone number or other recipient

    rendered_subject VARCHAR(255), -- rendered subject of the notification
    rendered_message TEXT, -- rendered message of the notification

    -- data JSONB,                          -- PER-USER override (optional)

    state VARCHAR(20) NOT NULL,           -- pending / processing / sending / failed / completed
    retry_count INT DEFAULT 0, -- number of times the notification has been retried
    error_message TEXT, -- error message from the last attempt

    sent_at TIMESTAMP, -- time the notification was sent
    created_at TIMESTAMP NOT NULL DEFAULT now(), -- notification log creation time

    CONSTRAINT fk_notification_logs
        FOREIGN KEY (notification_id)
        REFERENCES notifications(id)
);

CREATE INDEX idx_notification_logs_notification ON notification_logs(notification_id);
CREATE INDEX idx_notification_logs_user ON notification_logs(user_id);
CREATE INDEX idx_notification_logs_channel ON notification_logs(channel);
CREATE INDEX idx_notification_logs_state ON notification_logs(state);
CREATE INDEX idx_notification_logs_sent_at ON notification_logs(sent_at);
CREATE INDEX idx_notification_logs_created_at ON notification_logs(created_at);