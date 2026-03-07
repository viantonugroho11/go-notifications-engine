package kafka

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	libkafka "github.com/viantonugroho11/go-lib/kafka"
)

// Consumers menjalankan satu consumer group (satu topic + groupID). Bisa untuk tipe event apa saja (generic E).
// Lifecycle: Start(ctx) lalu Close().
type Consumers struct {
	c libkafka.Consumer
}

// Start menjalankan consumer sampai ctx dibatalkan.
func (c *Consumers) Start(ctx context.Context) {
	if c.c != nil {
		c.c.Start(ctx)
	}
}

// Close menutup consumer dan release resource.
func (c *Consumers) Close() error {
	if c.c == nil {
		return nil
	}
	return c.c.Close()
}

// NewConsumers membuat Consumers untuk tipe event E. Dipakai oleh RegisterConsumers[E].
func NewConsumers[E any](
	brokers []string,
	groupID string,
	topic string,
	handler libkafka.EventHandler[E],
	opts ...libkafka.ConsumerOption,
) (*Consumers, error) {
	c, err := libkafka.NewConsumer(brokers, groupID, topic, handler, opts...)
	if err != nil {
		return nil, err
	}
	return &Consumers{c: c}, nil
}

// RegisterConsumers membuat satu consumer untuk tipe event E dari config pertama.
// Handlers map[key]EventHandler[E]; key = consumer key (sama dengan config.HandlerKey).
// Kalau event message beda tipe, buat router/handler map tipe lain dan panggil RegisterConsumers[TipeLain].
func RegisterConsumers[E any](
	ctx context.Context,
	brokers []string,
	configs []ConsumerConfig,
	handlers map[string]libkafka.EventHandler[E],
	opts ...libkafka.ConsumerOption,
) (*Consumers, error) {
	if len(configs) == 0 {
		return nil, fmt.Errorf("consumer configs tidak boleh kosong")
	}
	cfg := configs[0]
	h, ok := handlers[cfg.HandlerKey]
	if !ok {
		return nil, fmt.Errorf("handler tidak ditemukan untuk key %q", cfg.HandlerKey)
	}
	consumerOpts := make([]libkafka.ConsumerOption, 0, len(opts)+2)
	if cfg.ClientID != "" {
		consumerOpts = append(consumerOpts, libkafka.WithConsumerClientID(cfg.ClientID))
	}
	consumerOpts = append(consumerOpts, libkafka.WithInitialOffset(sarama.OffsetOldest))
	consumerOpts = append(consumerOpts, opts...)
	return NewConsumers[E](brokers, cfg.GroupID, cfg.Topic, h, consumerOpts...)
}
