package notificationinbox

import (
	"context"

	inboxEntity "github.com/viantonugroho11/go-notifications-engine/internal/entity/notificationinbox"
)

type NotificationInboxRepository interface {
	Create(ctx context.Context, i inboxEntity.NotificationInbox) (inboxEntity.NotificationInbox, error)
	GetByID(ctx context.Context, id string) (inboxEntity.NotificationInbox, error)
	List(ctx context.Context, param *inboxEntity.NotificationInboxListParam) ([]inboxEntity.NotificationInbox, error)
	Update(ctx context.Context, i inboxEntity.NotificationInbox) (inboxEntity.NotificationInbox, error)
	Delete(ctx context.Context, id string) error
}
