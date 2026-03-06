package bootstrap

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	notifEntity "go-boilerplate-clean/internal/entity/notifications"
	"go-boilerplate-clean/internal/config"
	notifpg "go-boilerplate-clean/internal/repository/notification/postgres"
	tplpg "go-boilerplate-clean/internal/repository/notificationtemplate/postgres"
	eventhandler "go-boilerplate-clean/internal/transport/event/kafka/handler"
	eventkafka "go-boilerplate-clean/internal/transport/event/kafka"
	usecasenotif "go-boilerplate-clean/internal/usecase/notifications"

	"github.com/viantonugroho11/go-lib/kafka"
)

// ConsumerApp menjalankan Kafka consumer(s) dan mengelola lifecycle (DB, consumers, shutdown).
type ConsumerApp struct {
	cfg      config.Configuration
	consumers *eventkafka.Consumers
	close    func()
}

// NewConsumer membuat ConsumerApp: DB, migrate, repos, usecase, handler config, register consumers.
// Pemanggil wajib defer app.Close() lalu app.Run().
func NewConsumer(cfg config.Configuration) (*ConsumerApp, error) {
	ctx := context.Background()

	db, err := newDB(ctx, cfg)
	if err != nil {
		return nil, err
	}

	notificationRepo := notifpg.NewNotificationRepository(db)
	templateRepo := tplpg.NewNotificationTemplateRepository(db)
	notificationService := usecasenotif.NewNotificationService(notificationRepo, templateRepo, nil)

	configs := consumerConfigs(cfg)
	handlers := consumerHandlers(notificationService, nil) // sentSender nil
	consumers, err := eventhandler.RegisterConsumers(ctx, cfg.KafkaBrokersList(), configs, handlers)
	if err != nil {
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
		return nil, err
	}

	sqlDB, _ := db.DB()
	close := func() {
		_ = consumers.Close()
		_ = sqlDB.Close()
	}

	return &ConsumerApp{
		cfg:      cfg,
		consumers: consumers,
		close:    close,
	}, nil
}

func consumerConfigs(cfg config.Configuration) []eventhandler.ConsumerConfig {
	configs := []eventhandler.ConsumerConfig{
		{Topic: cfg.Kafka.Topic, GroupID: cfg.Kafka.GroupID, HandlerKey: "notification"},
	}
	if cfg.Kafka.TopicSent != "" {
		configs = append(configs, eventhandler.ConsumerConfig{
			Topic: cfg.Kafka.TopicSent, GroupID: cfg.Kafka.GroupID + "-sent", HandlerKey: "sent",
		})
	}
	return configs
}

func consumerHandlers(notificationService usecasenotif.NotificationService, sentSender eventhandler.NotificationSender) map[string]kafka.EventHandler[notifEntity.NotificationProducerMessage] {
	return map[string]kafka.EventHandler[notifEntity.NotificationProducerMessage]{
		"notification": eventhandler.NewNotificationUpdateHandler(notificationService),
		"sent":         eventhandler.NewNotificationSentHandler(sentSender),
	}
}

// Run block sampai dapat signal interrupt.
func (c *ConsumerApp) Run() {
	log.Printf("consumer started, group=%s topic=%s", c.cfg.Kafka.GroupID, c.cfg.Kafka.Topic)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("shutdown signal received, stopping consumer...")
}

// Close menutup consumers dan DB.
func (c *ConsumerApp) Close() {
	if c.close != nil {
		c.close()
	}
}
