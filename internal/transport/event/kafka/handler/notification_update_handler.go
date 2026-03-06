package handler

import (
	"context"
	"log"

	notifEntity "go-boilerplate-clean/internal/entity/notifications"

	kafka "github.com/viantonugroho11/go-lib/kafka"
)

// NotificationUpdater dependency untuk handler update notification dari event Kafka.
type NotificationUpdater interface {
	Update(ctx context.Context, n notifEntity.Notification) (notifEntity.Notification, error)
}

// NotificationUpdateHandler menangani event Kafka untuk update notification.
// Mengimplementasikan go-lib EventHandler[NotificationProducerMessage]; decode JSON dilakukan oleh consumer.
type NotificationUpdateHandler struct {
	updater NotificationUpdater
}

// NewNotificationUpdateHandler membuat handler update notification (inject dependency di sini).
func NewNotificationUpdateHandler(updater NotificationUpdater) *NotificationUpdateHandler {
	return &NotificationUpdateHandler{updater: updater}
}

// Name untuk logging/metrics consumer go-lib.
func (h *NotificationUpdateHandler) Name() string { return "notification-update" }

// Handle memproses event NotificationProducerMessage lalu panggil usecase Update.
func (h *NotificationUpdateHandler) Handle(ctx context.Context, evt notifEntity.NotificationProducerMessage, _ ...kafka.Header) kafka.Progress {
	n := evt.ToNotification()
	_, err := h.updater.Update(ctx, n)
	if err != nil {
		log.Printf("kafka notification update: update error: %v", err)
		return kafka.Progress{Status: kafka.ProgressError, Err: err}
	}
	return kafka.Progress{Status: kafka.ProgressSuccess}
}
