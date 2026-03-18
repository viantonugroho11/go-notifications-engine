package notifications

import (
	"context"

	notifEntity "go-boilerplate-clean/internal/entity/notifications"
)

// NotificationEventUsecase usecase untuk event notifikasi (e.g. consumer Kafka).
type NotificationEventUsecase interface {
	SendNotification(ctx context.Context, n notifEntity.Notification) error
}
