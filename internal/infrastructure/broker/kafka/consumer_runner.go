package kafka

import (
	"context"
	"fmt"

	"github.com/viantonugroho11/go-notifications-engine/internal/config"
	notifEntity "github.com/viantonugroho11/go-notifications-engine/internal/entity/notifications"
	eventkafka "github.com/viantonugroho11/go-notifications-engine/internal/transport/event/kafka"

	libkafka "github.com/viantonugroho11/go-lib/kafka"
)

var ErrConsumerKeyRequired = fmt.Errorf("flag -consumer wajib diisi")

func ErrConsumerUnknown(key string) error {
	return fmt.Errorf("consumer tidak dikenal: %q (tersedia: %v)", key, eventkafka.AllKeys)
}

// NewConsumerRunnerFor membuat consumer untuk tipe event E. Handlers = map[key]EventHandler[E].
// Pakai ini kalau event message beda tipe: buat router (HandlersX) di event/kafka, daftarkan key+config di registry, lalu panggil ini.
func NewConsumerRunnerFor[E any](
	cfg config.Configuration,
	consumerKey string,
	handlers map[string]libkafka.EventHandler[E],
) (*Consumers, error) {
	if consumerKey == "" {
		return nil, ErrConsumerKeyRequired
	}
	if !ValidConsumerKey(consumerKey) {
		return nil, ErrConsumerUnknown(consumerKey)
	}
	configToRun, err := GetConsumerConfigByKey(cfg, consumerKey)
	if err != nil {
		return nil, err
	}
	return RegisterConsumers[E](context.Background(), cfg.KafkaBrokersList(), []ConsumerConfig{configToRun}, handlers)
}

// NewConsumerRunner consumer untuk event notifikasi (NotificationsEventMessage: Action, After, Before). Shortcut yang pakai event/kafka.Handlers.
func NewConsumerRunner(cfg config.Configuration, consumerKey string, svc eventkafka.EventServices) (*Consumers, error) {
	return NewConsumerRunnerFor[notifEntity.NotificationsEventMessage](cfg, consumerKey, eventkafka.Handlers(svc))
}
