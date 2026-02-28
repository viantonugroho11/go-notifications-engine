package notification

import (
	"context"

	notifEntity "go-boilerplate-clean/internal/entity/notifications"
)

type NotificationRepository interface {
	Create(ctx context.Context, n notifEntity.Notification) (notifEntity.Notification, error)
	GetByID(ctx context.Context, id string) (notifEntity.Notification, error)
	List(ctx context.Context) ([]notifEntity.Notification, error)
	Update(ctx context.Context, n notifEntity.Notification) (notifEntity.Notification, error)
	Delete(ctx context.Context, id string) error
}
