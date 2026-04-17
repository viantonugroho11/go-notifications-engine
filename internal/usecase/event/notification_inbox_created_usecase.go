package event

import (
	"context"

	"github.com/viantonugroho11/go-notifications-engine/internal/client/notification"
	"github.com/viantonugroho11/go-notifications-engine/internal/entity/notificationinbox"
	"github.com/viantonugroho11/go-notifications-engine/internal/entity/notifications"
)

type NotificationInboxCreatedUsecase interface {
	Create(ctx context.Context, n notifications.NotificationEventUsecase) error
}

type notificationInboxCreatedUsecase struct {
	notificationlogClient notification.Client
}

func NewNotificationInboxCreatedUsecase(notificationlogClient notification.Client) NotificationInboxCreatedUsecase {
	return &notificationInboxCreatedUsecase{notificationlogClient: notificationlogClient}
}

func (s *notificationInboxCreatedUsecase) Create(ctx context.Context, n notifications.NotificationEventUsecase) error {
	_, err := s.notificationlogClient.CreateInbox(ctx, notificationinbox.NotificationInbox{
		UserID:            n.NotificationLogs.UserID,
		NotificationLogID: n.NotificationLogs.ID,
		Message:           n.NotificationLogs.RenderedMessage,
		Subject:           n.NotificationLogs.RenderedSubject,
	})
	if err != nil {
		return err
	}
	return nil
}
