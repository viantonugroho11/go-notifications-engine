package notifications

import (
	"context"

	notifEntity "github.com/viantonugroho11/go-notifications-engine/internal/entity/notifications"
)

// NotificationEventUsecase usecase untuk event notifikasi (e.g. consumer Kafka).
type NotificationEventUsecase interface {
	SendNotification(ctx context.Context, n notifEntity.Notification) error
}
