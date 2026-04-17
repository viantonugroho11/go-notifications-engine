package notification

import (
	"context"

	notifEntity "github.com/viantonugroho11/go-notifications-engine/internal/entity/notifications"
)

type NotificationRepository interface {
	Create(ctx context.Context, n notifEntity.Notification) (notifEntity.Notification, error)
	GetByID(ctx context.Context, id string) (notifEntity.Notification, error)
	List(ctx context.Context, param *notifEntity.NotificationListParam) ([]notifEntity.Notification, error)
	Update(ctx context.Context, n notifEntity.Notification) (notifEntity.Notification, error)
	Delete(ctx context.Context, id string) error
}
