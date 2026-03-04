package firebase

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
)

type FirebaseClient interface {
	Send(ctx context.Context, message *messaging.Message) (messageID string, err error)
}

type firebaseClient struct {
	app *firebase.App
}

func NewFirebaseClient(app *firebase.App) FirebaseClient {
	return &firebaseClient{app: app}
}

func (c *firebaseClient) Send(ctx context.Context, message *messaging.Message) (string, error) {
	client, err := c.app.Messaging(ctx)
	if err != nil {
		return "", err
	}
	return client.Send(ctx, message)
}