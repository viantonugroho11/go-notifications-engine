package firebase

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
)

type FirebaseClient interface {
	Send(ctx context.Context, message Message) (messageID string, err error)
}

type firebaseClient struct {
	app *firebase.App
}

func NewFirebaseClient(app *firebase.App) FirebaseClient {
	return &firebaseClient{app: app}
}

func (c *firebaseClient) Send(ctx context.Context, message Message) (string, error) {
	client, err := c.app.Messaging(ctx)
	if err != nil {
		return "", err
	}
	return client.Send(ctx, &messaging.Message{
		Token: message.Token,
		Notification: &messaging.Notification{
			Title: message.Notification.Title,
			Body: message.Notification.Body,
			ImageURL: message.Notification.Image,
		},
		Android: &messaging.AndroidConfig{
			Priority: message.Android.Priority,
			Notification: &messaging.AndroidNotification{
				Title: message.Notification.Title,
				Body: message.Notification.Body,
				ImageURL: message.Notification.Image,
			},
		},
		Webpush: &messaging.WebpushConfig{
			Headers: message.Webpush.Headers,
			Data: message.Webpush.Data,
			Notification: &messaging.WebpushNotification{
				Title: message.Notification.Title,
				Body: message.Notification.Body,
				Image: message.Notification.Image,
			},
		},
		APNS: &messaging.APNSConfig{
			Headers: message.APNS.Headers,
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Alert: &messaging.ApsAlert{
						Title: message.Notification.Title,
						Body: message.Notification.Body,
						LaunchImage: message.Notification.Image,
					},
				},
			},
		},
		FCMOptions: &messaging.FCMOptions{
			AnalyticsLabel: message.FCMOptions.AnalyticsLabel,
		},
		Topic: message.Topic,
		Condition: message.Condition,
		Data: message.Data,
	})
}