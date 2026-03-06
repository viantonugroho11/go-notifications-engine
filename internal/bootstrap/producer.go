package bootstrap

import (
	"context"

	"go-boilerplate-clean/internal/config"
	notifEntity "go-boilerplate-clean/internal/entity/notifications"
	usecasenotif "go-boilerplate-clean/internal/usecase/notifications"

	"github.com/IBM/sarama"
	"github.com/viantonugroho11/go-lib/kafka"
)

// kafkaProducerCloser mempublish event notifikasi dan bisa ditutup (untuk cleanup).
type kafkaProducerCloser interface {
	usecasenotif.NotificationEventPublisher
	Close() error
}

// newNotificationProducer membuat producer go-lib untuk NotificationProducerMessage.
func newNotificationProducer(cfg config.Configuration) (kafkaProducerCloser, error) {
	p, err := kafka.NewProducer[notifEntity.NotificationProducerMessage](
		cfg.KafkaBrokersList(),
		cfg.Kafka.Topic,
		kafka.WithKeyFunc(func(m notifEntity.NotificationProducerMessage) []byte { return []byte(m.NotificationID) }),
		kafka.WithAcks(sarama.WaitForAll),
		kafka.WithIdempotent(),
		kafka.WithRetryMax(5),
	)
	if err != nil {
		return nil, err
	}
	return &notificationProducerAdapter{p: p}, nil
}

// notificationProducerAdapter mengadapt go-lib Producer ke usecase NotificationEventPublisher.
type notificationProducerAdapter struct {
	p *kafka.Producer[notifEntity.NotificationProducerMessage]
}

func (a *notificationProducerAdapter) Publish(ctx context.Context, msg notifEntity.NotificationProducerMessage) error {
	return a.p.Publish(ctx, msg)
}

func (a *notificationProducerAdapter) Close() error {
	return a.p.Close()
}
