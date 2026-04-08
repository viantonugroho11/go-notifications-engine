package event

import (
	"go-boilerplate-clean/internal/client/notification"
)

type NotificationInboxCreatedUsecase interface {
	// Create(ctx context.Context, n Notification) error
}

type notificationInboxCreatedUsecase struct {
	notificationlogClient notification.Client
}

// func NewNotificationInboxCreatedUsecase(notificationlogClient notificationlogClient) NotificationInboxCreatedUsecase {
// 	return &notificationInboxCreatedUsecase{notificationlogClient: notificationlogClient}
// }

// func (s *notificationInboxCreatedUsecase) Create(ctx context.Context, n Notification) error {
// 	return s.notificationlogClient.Create(ctx, n)
// }