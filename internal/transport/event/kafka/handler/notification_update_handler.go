package handler

import (
	"context"
	"log"
	"strings"

	notifEntity "go-boilerplate-clean/internal/entity/notifications"

	kafka "github.com/viantonugroho11/go-lib/kafka"
)

// NotificationUpdater dependency untuk handler update notification dari event Kafka.
type NotificationUpdater interface {
	Update(ctx context.Context, n notifEntity.Notification) (notifEntity.Notification, error)
}

// NotificationUpdateHandler menangani event Kafka (NotificationsEventMessage: Action, After, Before).
type NotificationUpdateHandler struct {
	updater NotificationUpdater
}

// NewNotificationUpdateHandler membuat handler update notification (inject dependency di sini).
func NewNotificationUpdateHandler(updater NotificationUpdater) *NotificationUpdateHandler {
	return &NotificationUpdateHandler{updater: updater}
}

// Name untuk logging/metrics consumer go-lib.
func (h *NotificationUpdateHandler) Name() string { return "notification-update" }

// Handle memproses NotificationsEventMessage; untuk INSERT/UPDATE pakai After, untuk DELETE skip atau pakai Before sesuai kebutuhan.
func (h *NotificationUpdateHandler) Handle(ctx context.Context, evt notifEntity.NotificationsEventMessage, _ ...kafka.Header) kafka.Progress {
	switch strings.ToUpper(evt.Action) {
	case "DELETE":
		return kafka.Progress{Status: kafka.ProgressSkip}
	}

	if evt.After.State == notifEntity.StateScheduled.String() {
		return kafka.Progress{Status: kafka.ProgressSkip}
	}

	if evt.After.State == notifEntity.StateProcessing.String() {
		return kafka.Progress{Status: kafka.ProgressSkip}
	}



	n := evt.After.ToNotification()
	_, err := h.updater.Update(ctx, n)
	if err != nil {
		log.Printf("kafka notification update: update error: %v", err)
		return kafka.Progress{Status: kafka.ProgressError, Err: err}
	}
	return kafka.Progress{Status: kafka.ProgressSuccess}
}
