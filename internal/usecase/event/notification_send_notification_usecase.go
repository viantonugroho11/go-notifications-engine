package event

import (
	"context"
	"errors"
	"time"

	"github.com/viantonugroho11/go-notifications-engine/internal/client/email"
	"github.com/viantonugroho11/go-notifications-engine/internal/client/firebase"
	"github.com/viantonugroho11/go-notifications-engine/internal/client/notification"
	"github.com/viantonugroho11/go-notifications-engine/internal/entity/notificationlogs"
	"github.com/viantonugroho11/go-notifications-engine/internal/entity/notifications"
)

type NotificationProcessingToUsecase interface {
	Process(ctx context.Context, n notifications.NotificationEventUsecase) error
}

type notificationProcessingToUsecase struct {
	notificationClient notification.Client
	emailClient        email.EmailClient
	firebaseClient     firebase.FirebaseClient
}

func NewNotificationProcessingToUsecase(notificationClient notification.Client, emailClient email.EmailClient, firebaseClient firebase.FirebaseClient) NotificationProcessingToUsecase {
	return &notificationProcessingToUsecase{notificationClient: notificationClient, emailClient: emailClient, firebaseClient: firebaseClient}
}

func (s *notificationProcessingToUsecase) Process(ctx context.Context, n notifications.NotificationEventUsecase) error {
	var remark string
	var status notificationlogs.State

	switch n.Channel {
	case notifications.ChannelEmail:
		err := s.emailClient.Send(ctx, &email.Message{
			To:      []string{n.NotificationLogs.SendTo},
			Subject: n.NotificationLogs.RenderedSubject,
			Body:    n.NotificationLogs.RenderedMessage,
		})
		if err != nil {
			remark = err.Error()
		}
	case notifications.ChannelPush:
		messageID, err := s.firebaseClient.Send(ctx, firebase.Message{
			Token: n.NotificationLogs.SendTo,
			Notification: &firebase.Notification{
				Title: n.NotificationLogs.RenderedSubject,
				Body:  n.NotificationLogs.RenderedMessage,
			},
		})
		if err != nil {
			remark = err.Error()
		}
		remark = messageID
	default:
		return errors.New("channel not supported")
	}

	status = notificationlogs.StateCompleted
	if remark != "" {
		status = notificationlogs.StateFailed
	}

	now := time.Now()
	_, err := s.notificationClient.UpdateNotificationLog(ctx, notificationlogs.NotificationLog{
		ID:           n.NotificationLogs.ID,
		State:        status,
		ErrorMessage: remark,
		SentAt:       &now,
	})
	if err != nil {
		return err
	}
	return nil
}
