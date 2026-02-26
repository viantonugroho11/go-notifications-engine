CREATE TABLE notification_logs (
    id UUID PRIMARY KEY,
    notification_id UUID NOT NULL,
    user_id UUID NOT NULL,

    channel VARCHAR(50) NOT NULL,
    send_to TEXT,

    rendered_subject VARCHAR(255),
    rendered_message TEXT,

    data JSONB,                          -- PER-USER override (optional)

    state VARCHAR(20) NOT NULL,           -- queued / processing / sent / failed / completed
    retry_count INT DEFAULT 0,
    error_message TEXT,

    sent_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT now(),

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