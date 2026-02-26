CREATE TABLE notification_inbox (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    notification_log_id UUID NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT now(),

    CONSTRAINT fk_notification_inbox_log
        FOREIGN KEY (notification_log_id)
        REFERENCES notification_logs(id),

    CONSTRAINT uq_notification_user_log UNIQUE (notification_log_id)
);

CREATE INDEX idx_notification_inbox_user ON notification_inbox(user_id);
CREATE INDEX idx_notification_inbox_log ON notification_inbox(notification_log_id);
CREATE INDEX idx_notification_inbox_read ON notification_inbox(is_read);
CREATE INDEX idx_notification_inbox_read_at ON notification_inbox(read_at);
CREATE INDEX idx_notification_inbox_created_at ON notification_inbox(created_at);