package bootstrap

import (
	"github.com/viantonugroho11/go-notifications-engine/internal/config"
	brokerkafka "github.com/viantonugroho11/go-notifications-engine/internal/infrastructure/broker/kafka"
)

// newNotificationProducer membuat producer notifikasi via broker (go-lib/kafka).
func newNotificationProducer(cfg config.Configuration) (brokerkafka.NotificationProducerCloser, error) {
	return brokerkafka.NewNotificationProducer(cfg)
}
