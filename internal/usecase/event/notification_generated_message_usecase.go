package event

import (
	"context"
	"go-boilerplate-clean/internal/client/notification"
	"go-boilerplate-clean/internal/entity/notificationlogs"
	"go-boilerplate-clean/internal/entity/notifications"
)

type NotificationGeneratedMessageUsecase interface {
	GenerateMessage(ctx context.Context, n notifications.NotificationEventUsecase) error
}

type notificationGeneratedMessageUsecase struct {
	notificationClient notification.Client
}

func NewNotificationGeneratedMessageUsecase(notificationClient notification.Client) NotificationGeneratedMessageUsecase {
	return &notificationGeneratedMessageUsecase{notificationClient: notificationClient}
}

func (s *notificationGeneratedMessageUsecase) GenerateMessage(ctx context.Context, n notifications.NotificationEventUsecase) error {
	
	notificationTemplate, err := s.notificationClient.GetNotificationTemplate(ctx, n.NotificationTemplateID)
	if err != nil {
		return err
	}
	renderedMessage := n.GenerateRenderedMessage(notificationTemplate.Body)
	renderedSubject := n.GenerateRenderedSubject(notificationTemplate.Subject)
	n.NotificationLogs.RenderedMessage = renderedMessage
	n.NotificationLogs.RenderedSubject = renderedSubject
	_, err = s.notificationClient.UpdateNotificationLog(ctx, notificationlogs.NotificationLog{
		ID: n.NotificationLogs.ID,
		RenderedMessage: renderedMessage,
		RenderedSubject: renderedSubject,
		State: notificationlogs.StateSent,
	})
	if err != nil {
		return err
	}			
	return nil
}