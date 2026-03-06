package handler

import (
	"context"
	"fmt"

	notifEntity "go-boilerplate-clean/internal/entity/notifications"
	eventkafka "go-boilerplate-clean/internal/transport/event/kafka"

	"github.com/IBM/sarama"
	kafka "github.com/viantonugroho11/go-lib/kafka"
)

// ConsumerConfig konfigurasi satu consumer: topic, group, dan key handler yang dipakai.
type ConsumerConfig struct {
	Topic      string
	GroupID    string
	HandlerKey string
}

// RegisterConsumers membuat satu consumer per config menggunakan go-lib/kafka.
// Handlers harus implement kafka.EventHandler[NotificationProducerMessage]; didaftarkan per key dari main.
// Pemanggil wajib defer consumers.Close().
func RegisterConsumers(ctx context.Context, brokers []string, configs []ConsumerConfig, handlers map[string]kafka.EventHandler[notifEntity.NotificationProducerMessage]) (*eventkafka.Consumers, error) {
	out := &eventkafka.Consumers{}
	for _, cfg := range configs {
		h := handlers[cfg.HandlerKey]
		if h == nil {
			return nil, fmt.Errorf("event/handler: handler tidak terdaftar untuk key %q", cfg.HandlerKey)
		}
		consumer, err := kafka.NewConsumer[notifEntity.NotificationProducerMessage](
			brokers,
			cfg.GroupID,
			cfg.Topic,
			h,
			kafka.WithInitialOffset(sarama.OffsetOldest),
		)
		if err != nil {
			return nil, err
		}
		out.List = append(out.List, consumer)
		go consumer.Start(ctx)
	}
	return out, nil
}
