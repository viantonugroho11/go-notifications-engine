package notification

import (
	"context"
	"log"

	"github.com/viantonugroho11/go-notifications-engine/internal/client/email"
	"github.com/viantonugroho11/go-notifications-engine/internal/client/firebase"
	notifEntity "github.com/viantonugroho11/go-notifications-engine/internal/entity/notifications"
)

// SentSender mengirim notifikasi via email atau Firebase berdasarkan channel.
type SentSender struct {
	Email    email.EmailClient
	Firebase firebase.FirebaseClient
}

// Send mengirim sesuai channel: email -> EmailClient, push -> FirebaseClient.
func (s *SentSender) Send(ctx context.Context, msg *notifEntity.NotificationProducerMessage) error {
	if msg == nil {
		return nil
	}
	ch := msg.Channel
	for i := range msg.NotificationLogs {
		logItem := &msg.NotificationLogs[i]
		sendTo := logItem.SendTo
		subject := logItem.RenderedSubject
		body := logItem.RenderedMessage
		switch ch {
		case "email":
			if s.Email != nil && sendTo != "" {
				if err := s.Email.Send(ctx, &email.Message{
					To:      []string{sendTo},
					Subject: subject,
					Body:    body,
				}); err != nil {
					log.Printf("notification sender: email error: %v", err)
					return err
				}
			}
		case "push":
			if s.Firebase != nil && sendTo != "" {
				if _, err := s.Firebase.Send(ctx, firebase.Message{
					Token: sendTo,
					Notification: &firebase.Notification{
						Title: subject,
						Body:  body,
					},
				}); err != nil {
					log.Printf("notification sender: firebase error: %v", err)
					return err
				}
			}
		default:
			log.Printf("notification sender: channel %q not implemented", ch)
		}
	}
	return nil
}
