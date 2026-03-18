CREATE TABLE notifications (
    id UUID PRIMARY KEY,
    event_key VARCHAR(100) NOT NULL,
    notification_template_id UUID NOT NULL,
    data JSONB,
    category VARCHAR(20) NOT NULL, -- promo / transactional / system / other
    channel VARCHAR(20) NOT NULL, -- email / sms / push / whatsapp / telegram / line / wechat / weibo / kakao
    
    state VARCHAR(20) NOT NULL, -- created / scheduled / processing / sent / failed / completed

    schedule_at TIMESTAMP, -- schedule time for sending notification
    created_by VARCHAR(100) DEFAULT 'system', -- user who created the notification
    updated_by VARCHAR(100), -- user who updated the notification
    created_at TIMESTAMP NOT NULL DEFAULT now(), -- notification creation time
    updated_at TIMESTAMP, -- notification update time

    CONSTRAINT fk_notifications_template
        FOREIGN KEY (notification_template_id)
        REFERENCES notification_templates(id),

    CONSTRAINT uq_event_key UNIQUE (event_key)
);


CREATE INDEX idx_notifications_event ON notifications(event_key);
CREATE INDEX idx_notifications_template ON notifications(notification_template_id);
CREATE INDEX idx_notifications_state ON notifications(state);
CREATE INDEX idx_notifications_send_time ON notifications(send_time);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);
CREATE INDEX idx_notifications_updated_at ON notifications(updated_at);