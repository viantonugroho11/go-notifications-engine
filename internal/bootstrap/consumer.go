package bootstrap

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/viantonugroho11/go-notifications-engine/internal/client/email"
	"github.com/viantonugroho11/go-notifications-engine/internal/client/firebase"
	clientnotif "github.com/viantonugroho11/go-notifications-engine/internal/client/notification"
	"github.com/viantonugroho11/go-notifications-engine/internal/config"
	brokerkafka "github.com/viantonugroho11/go-notifications-engine/internal/infrastructure/broker/kafka"
	notifpg "github.com/viantonugroho11/go-notifications-engine/internal/repository/notification/postgres"
	tplpg "github.com/viantonugroho11/go-notifications-engine/internal/repository/notificationtemplate/postgres"
	eventkafka "github.com/viantonugroho11/go-notifications-engine/internal/transport/event/kafka"
	usecasenotif "github.com/viantonugroho11/go-notifications-engine/internal/usecase/notifications"
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

	// HTTP client ke API notification (self-call): consumer mutasi state lewat API agar
	// business logic (state machine, validasi, publish event) tetap di satu tempat.
	// Timeout 10s: consumer adalah async, network error → Kafka retry.
	httpClient := &http.Client{Timeout: 10 * time.Second}
	_ = clientnotif.NewClient(cfg.App.NotificationBaseURL, httpClient) // dipakai oleh event usecases jika di-wire

	// Email client
	var emailClient email.EmailClient
	if cfg.Email.Host != "" {
		dialer, err := config.InitializeEmail(cfg.Email)
		if err != nil {
			log.Printf("consumer: email client init failed: %v (email tidak akan terkirim)", err)
		} else {
			emailClient = email.NewEmailClient(dialer, cfg.Email.User)
		}
	}

	// Firebase client
	var firebaseClient firebase.FirebaseClient
	if cfg.FCM.ProjectID != "" {
		fbApp, err := config.ConnectToFirebase(cfg.FCM.ProjectID)
		if err != nil {
			log.Printf("consumer: firebase client init failed: %v (push tidak akan terkirim)", err)
		} else {
			firebaseClient = firebase.NewFirebaseClient(fbApp)
		}
	}

	sentSender := &clientnotif.SentSender{
		Email:    emailClient,
		Firebase: firebaseClient,
	}

	svc := eventkafka.EventServices{Notification: notificationService, SentSender: sentSender}
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
