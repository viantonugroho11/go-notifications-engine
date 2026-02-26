CREATE TABLE notification_templates (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,          -- internal key
    subject VARCHAR(255),                -- push title / email subject
    body TEXT,                            -- push / sms body / email html
    payload_schema jsonb,                 -- payload schema for validation
    channel VARCHAR(50) NOT NULL,        -- push / email / sms
    template_type VARCHAR(50),           -- promo / transactional
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_notification_templates_channel ON notification_templates(channel);
CREATE INDEX idx_notification_templates_type ON notification_templates(template_type);
CREATE INDEX idx_notification_templates_deleted ON notification_templates(deleted_at);