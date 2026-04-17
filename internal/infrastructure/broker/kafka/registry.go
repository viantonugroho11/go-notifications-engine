package kafka

import (
	"fmt"

	"github.com/viantonugroho11/go-notifications-engine/internal/config"
	eventkafka "github.com/viantonugroho11/go-notifications-engine/internal/transport/event/kafka"
)

// ConsumerConfig konfigurasi per consumer (topic, group, clientID). HandlerKey = event key.
type ConsumerConfig struct {
	Topic      string
	GroupID    string
	ClientID   string
	HandlerKey string
}

// ValidConsumerKey true jika key terdaftar di event layer (semua tipe event).
func ValidConsumerKey(key string) bool {
	for _, k := range eventkafka.AllKeys {
		if k == key {
			return true
		}
	}
	return false
}

// AvailableConsumerKeys untuk flag/help (delegasi ke event).
func AvailableConsumerKeys() []string {
	return eventkafka.AvailableKeys()
}

// GetConsumerConfigByKey mengembalikan config Kafka untuk key (topic, group dari cfg).
// Key mengikuti event/kafka; tambah case untuk consumer baru.
func GetConsumerConfigByKey(cfg config.Configuration, key string) (ConsumerConfig, error) {
	switch key {
	case eventkafka.KeyNotification:
		return ConsumerConfig{
			Topic:      cfg.Kafka.Topic,
			GroupID:    cfg.Kafka.GroupID,
			ClientID:   cfg.Kafka.ClientID,
			HandlerKey: eventkafka.KeyNotification,
		}, nil
	case eventkafka.KeySent:
		if cfg.Kafka.TopicSent == "" {
			return ConsumerConfig{}, fmt.Errorf("consumer %q memerlukan konfigurasi topic_sent", key)
		}
		return ConsumerConfig{
			Topic:      cfg.Kafka.TopicSent,
			GroupID:    cfg.Kafka.GroupID + "-sent",
			ClientID:   cfg.Kafka.ClientID + "-sent",
			HandlerKey: eventkafka.KeySent,
		}, nil
	default:
		return ConsumerConfig{}, fmt.Errorf("consumer tidak dikenal: %q (tersedia: %v)", key, eventkafka.AllKeys)
	}
}
