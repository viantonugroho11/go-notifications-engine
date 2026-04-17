package event

import (
	"context"
	"errors"

	"github.com/viantonugroho11/go-notifications-engine/internal/client/notification"
	"github.com/viantonugroho11/go-notifications-engine/internal/client/person"
	"github.com/viantonugroho11/go-notifications-engine/internal/entity/notificationlogs"
	"github.com/viantonugroho11/go-notifications-engine/internal/entity/notifications"
)

type NotificationFetchPersonUsecase interface {
	Fetch(ctx context.Context, n notifications.NotificationEventUsecase) error
}

type notificationFetchPersonUsecase struct {
	personClient       person.PersonClient
	notificationClient notification.Client
}

func NewNotificationFetchPersonUsecase(personClient person.PersonClient, notificationClient notification.Client) NotificationFetchPersonUsecase {
	return &notificationFetchPersonUsecase{personClient: personClient, notificationClient: notificationClient}
}

func (s *notificationFetchPersonUsecase) Fetch(ctx context.Context, n notifications.NotificationEventUsecase) error {
	var sendTo string
	person, err := s.personClient.GetPerson(ctx, n.NotificationLogs.UserID)
	if err != nil {
		return err
	}
	switch n.Channel {
	case notifications.ChannelEmail:
		sendTo = person.Email
	case notifications.ChannelPush:
		// Find the firebase device with the most recent LastActiveAt
		var latestDeviceToken string
		var latestTime int64
		for _, device := range person.Devices {
			if device.LastActiveAt.UnixNano() > latestTime {
				latestTime = device.LastActiveAt.UnixNano()
			}
		}
		sendTo = latestDeviceToken
	case notifications.ChannelSMS:
		sendTo = person.Phone
	case notifications.ChannelWhatsApp:
		sendTo = person.Phone
	case notifications.ChannelTelegram:
		sendTo = person.Phone
	case notifications.ChannelLine:
		sendTo = person.Phone
	default:
		return errors.New("channel not supported")
	}
	n.NotificationLogs.SendTo = sendTo

	_, err = s.notificationClient.UpdateNotificationLog(ctx, notificationlogs.NotificationLog{
		ID:     n.NotificationLogs.ID,
		SendTo: sendTo,
	})
	if err != nil {
		return err
	}
	return nil
}
