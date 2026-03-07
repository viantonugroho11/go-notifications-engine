package bootstrap

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-boilerplate-clean/internal/config"
	brokerkafka "go-boilerplate-clean/internal/infrastructure/broker/kafka"
	notifpg "go-boilerplate-clean/internal/repository/notification/postgres"
	tplpg "go-boilerplate-clean/internal/repository/notificationtemplate/postgres"
	eventkafka "go-boilerplate-clean/internal/transport/event/kafka"
	usecasenotif "go-boilerplate-clean/internal/usecase/notifications"
)

// ConsumerApp menjalankan satu Kafka consumer dan mengelola lifecycle (DB, consumer, shutdown).
// Routing (flag mana yang jalan) dan init consumer ada di broker/kafka.
type ConsumerApp struct {
	consumers      *brokerkafka.Consumers
	runConsumerKey string
	close          func()
}

// NewConsumer membuat ConsumerApp: DB, repos, usecase, lalu consumer dari broker/kafka (registry + routing by flag).
// Pemanggil wajib defer app.Close() lalu app.Run().
func NewConsumer(cfg config.Configuration, consumerKey string) (*ConsumerApp, error) {
	db, err := newDB(context.Background(), cfg)
	if err != nil {
		return nil, err
	}

	notificationRepo := notifpg.NewNotificationRepository(db)
	templateRepo := tplpg.NewNotificationTemplateRepository(db)
	notificationService := usecasenotif.NewNotificationService(notificationRepo, templateRepo, nil)

	svc := eventkafka.EventServices{Notification: notificationService, SentSender: nil}
	consumers, err := brokerkafka.NewConsumerRunner(cfg, consumerKey, svc)
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
		consumers:      consumers,
		runConsumerKey: consumerKey,
		close:          close,
	}, nil
}

// Run menjalankan consumer (go-lib) sampai dapat signal interrupt, lalu shutdown.
func (c *ConsumerApp) Run() {
	log.Printf("consumer started (run=%s)", c.runConsumerKey)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go c.consumers.Start(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("shutdown signal received, stopping consumer...")
	cancel()
}

// Close menutup consumers dan DB.
func (c *ConsumerApp) Close() {
	if c.close != nil {
		c.close()
	}
}
