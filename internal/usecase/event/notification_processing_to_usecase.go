package event

import (
	"context"
)

type NotificationProcessingToUsecase interface {
	Process(ctx context.Context, n Notification) error
}

type notificationProcessingToUsecase struct {
	notificationRepository notificationRepository
}

func NewNotificationProcessingToUsecase(notificationRepository notificationRepository) NotificationProcessingToUsecase {
	return &notificationProcessingToUsecase{notificationRepository: notificationRepository}
}