package email

import (
	"context"
	"errors"

	"gopkg.in/gomail.v2"
)


type EmailClient interface {
	Send(ctx context.Context, msg *Message) error
}

type emailClient struct {
	dialer *gomail.Dialer
	from   string // default From jika kosong di Message
}

func NewEmailClient(dialer *gomail.Dialer, from string) EmailClient {
	return &emailClient{dialer: dialer, from: from}
}

func (c *emailClient) Send(ctx context.Context, msg *Message) error {
	if msg == nil {
		return errors.New("email message is required")
	}
	if len(msg.To) == 0 {
		return errors.New("email recipient To is required")
	}
	m := gomail.NewMessage()
	from := msg.From
	if from == "" {
		from = c.from
	}
	m.SetHeader("From", from)
	m.SetHeader("To", msg.To...)
	m.SetHeader("Subject", msg.Subject)
	m.SetBody("text/plain", msg.Body)
	return c.dialer.DialAndSend(m)
}