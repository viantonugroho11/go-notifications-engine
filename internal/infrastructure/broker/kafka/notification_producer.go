package kafka

import (
	"context"

	"go-boilerplate-clean/internal/config"
	notifEntity "go-boilerplate-clean/internal/entity/notifications"
	usecasenotif "go-boilerplate-clean/internal/usecase/notifications"

	"github.com/IBM/sarama"
	libkafka "github.com/viantonugroho11/go-lib/kafka"
)

// NotificationProducerCloser mempublish event notifikasi dan bisa ditutup (untuk cleanup).
// Menggunakan go-lib/kafka Producer; implementasi dipakai oleh bootstrap.
type NotificationProducerCloser interface {
	usecasenotif.NotificationEventPublisher
	Close() error
}

// NewNotificationProducer membuat producer notifikasi dengan go-lib (satu topic, typed NotificationProducerMessage).
// Brokers diambil dari cfg.Kafka; topic dari cfg.Kafka.Topic.
func NewNotificationProducer(cfg config.Configuration) (NotificationProducerCloser, error) {
	brokers := cfg.KafkaBrokersList()
	p, err := libkafka.NewProducer[notifEntity.NotificationProducerMessage](
		brokers,
		cfg.Kafka.Topic,
		libkafka.WithKeyFunc(func(m notifEntity.NotificationProducerMessage) []byte {
			return []byte(m.NotificationID)
		}),
		libkafka.WithAcks(sarama.WaitForAll),
		libkafka.WithIdempotent(),
		libkafka.WithRetryMax(5),
	)
	if err != nil {
		return nil, err
	}
	return &notificationProducerAdapter{p: p}, nil
}

type notificationProducerAdapter struct {
	p *libkafka.Producer[notifEntity.NotificationProducerMessage]
}

func (a *notificationProducerAdapter) Publish(ctx context.Context, msg notifEntity.NotificationProducerMessage) error {
	return a.p.Publish(ctx, msg)
}

func (a *notificationProducerAdapter) Close() error {
	return a.p.Close()
}
