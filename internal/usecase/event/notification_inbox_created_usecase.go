package event

import (
	"context"
	"go-boilerplate-clean/internal/client/notification"
	"go-boilerplate-clean/internal/entity/notificationinbox"
	"go-boilerplate-clean/internal/entity/notifications"
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
		UserID: n.NotificationLogs.UserID,
		NotificationLogID: n.NotificationLogs.ID,
		Message: n.NotificationLogs.RenderedMessage,
		Subject: n.NotificationLogs.RenderedSubject,
	})
	if err != nil {
		return err
	}
	return nil
}