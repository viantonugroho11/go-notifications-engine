package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-boilerplate-clean/internal/config"
	notifEntity "go-boilerplate-clean/internal/entity/notifications"
	pginfra "go-boilerplate-clean/internal/infrastructure/database/postgres"
	notifpg "go-boilerplate-clean/internal/repository/notification/postgres"
	tplpg "go-boilerplate-clean/internal/repository/notificationtemplate/postgres"
	eventhandler "go-boilerplate-clean/internal/transport/event/kafka/handler"
	usecasenotif "go-boilerplate-clean/internal/usecase/notifications"

	confLoader "github.com/viantonugroho11/go-config-library"
	kafka "github.com/viantonugroho11/go-lib/kafka"
)

func main() {
	// Load configuration (Consul/env/file)
	cfg := config.Configuration{}
	loader := confLoader.New("", "go-boilerplate-clean", os.Getenv("CONSUL_URL"),
		confLoader.WithConfigFileSearchPaths("./config"),
	)
	if err := loader.Load(&cfg); err != nil {
		log.Fatalf("config load error: %v", err)
	}

	// Root context with cancellation on signal
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// DB untuk notification update dari event (consumer tidak publish, producer=nil)
	db, err := pginfra.Connect(ctx, cfg.PGDSN())
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err := pginfra.Migrate(db); err != nil {
		log.Fatalf("db migrate error: %v", err)
	}
	notificationRepo := notifpg.NewNotificationRepository(db)
	templateRepo := tplpg.NewNotificationTemplateRepository(db)
	notificationService := usecasenotif.NewNotificationService(notificationRepo, templateRepo, nil, "")

	// Handler didaftarkan per key; consumer memakai go-lib/kafka EventHandler[NotificationProducerMessage]
	consumerConfigs := []eventhandler.ConsumerConfig{
		{Topic: cfg.Kafka.Topic, GroupID: cfg.Kafka.GroupID, HandlerKey: "notification"},
	}
	if cfg.Kafka.TopicSent != "" {
		consumerConfigs = append(consumerConfigs, eventhandler.ConsumerConfig{
			Topic: cfg.Kafka.TopicSent, GroupID: cfg.Kafka.GroupID + "-sent", HandlerKey: "sent",
		})
	}
	var sentSender eventhandler.NotificationSender
	kafkaHandlers := map[string]kafka.EventHandler[notifEntity.NotificationProducerMessage]{
		"notification": eventhandler.NewNotificationUpdateHandler(notificationService),
		"sent":         eventhandler.NewNotificationSentHandler(sentSender),
	}
	consumers, err := eventhandler.RegisterConsumers(ctx, cfg.KafkaBrokersList(), consumerConfigs, kafkaHandlers)
	if err != nil {
		log.Fatalf("kafka consumer init error: %v", err)
	}
	defer func() {
		if cerr := consumers.Close(); cerr != nil {
			log.Printf("kafka consumers close error: %v", cerr)
		}
	}()
	log.Printf("consumer started, group=%s topic=%s", cfg.Kafka.GroupID, cfg.Kafka.Topic)

	// Wait for termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	log.Println("shutdown signal received, stopping consumer...")
}
