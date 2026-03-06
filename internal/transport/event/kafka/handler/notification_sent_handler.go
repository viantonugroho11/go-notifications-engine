package handler

import (
	"context"
	"log"

	notifEntity "go-boilerplate-clean/internal/entity/notifications"

	kafka "github.com/viantonugroho11/go-lib/kafka"
)

// NotificationSender mengirim notifikasi ke client (email/firebase) saat consume topic sent.
type NotificationSender interface {
	Send(ctx context.Context, msg *notifEntity.NotificationProducerMessage) error
}

// NotificationSentHandler consume topic sent lalu kirim via email/firebase.
// Mengimplementasikan go-lib EventHandler[NotificationProducerMessage].
type NotificationSentHandler struct {
	sender NotificationSender
}

// NewNotificationSentHandler membuat handler untuk topic sent.
func NewNotificationSentHandler(sender NotificationSender) *NotificationSentHandler {
	return &NotificationSentHandler{sender: sender}
}

// Name untuk logging/metrics consumer go-lib.
func (h *NotificationSentHandler) Name() string { return "notification-sent" }

// Handle memproses event lalu panggil sender (email/firebase).
func (h *NotificationSentHandler) Handle(ctx context.Context, evt notifEntity.NotificationProducerMessage, _ ...kafka.Header) kafka.Progress {
	if h.sender == nil {
		return kafka.Progress{Status: kafka.ProgressSkip}
	}
	if err := h.sender.Send(ctx, &evt); err != nil {
		log.Printf("kafka notification sent: send error: %v", err)
		return kafka.Progress{Status: kafka.ProgressError, Err: err}
	}
	return kafka.Progress{Status: kafka.ProgressSuccess}
}
