package bootstrap

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/viantonugroho11/go-notifications-engine/internal/config"
	redisinfra "github.com/viantonugroho11/go-notifications-engine/internal/infrastructure/cache/redis"
	"github.com/viantonugroho11/go-notifications-engine/internal/transport/apis"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// App menjalankan HTTP server dan mengelola lifecycle (DB, Kafka, Redis, shutdown).
type App struct {
	cfg    config.Configuration
	e      *echo.Echo
	server *http.Server
	close  func()
}

// New membuat App: koneksi DB, migrate, Kafka producer, Redis, repos, usecases, routes.
// Pemanggil wajib defer app.Close() lalu app.Run().
func New(cfg config.Configuration) (*App, error) {
	ctx := context.Background()

	db, err := newDB(ctx, cfg)
	if err != nil {
		return nil, err
	}

	publisher, err := newNotificationProducer(cfg)
	if err != nil {
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
		return nil, err
	}

	redisClient, err := redisinfra.NewClient(cfg.Redis.Addr, cfg.Redis.Password, strconv.Itoa(cfg.Redis.DB))
	if err != nil {
		_ = publisher.Close()
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
		return nil, err
	}

	svc := newServices(db, publisher)
	e := newEcho(svc)

	sqlDB, _ := db.DB()
	close := func() {
		_ = redisClient.Close()
		_ = publisher.Close()
		_ = sqlDB.Close()
	}

	return &App{
		cfg: cfg,
		e:   e,
		server: &http.Server{
			Addr:         ":" + cfg.App.Port,
			Handler:      e,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		close: close,
	}, nil
}

func newEcho(svc apis.Services) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover(), middleware.Logger())
	apis.RegisterRoutes(e, svc)
	return e
}

// Run menjalankan server, block sampai dapat signal interrupt, lalu shutdown.
func (a *App) Run() {
	go func() {
		if err := a.e.StartServer(a.server); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()
	log.Printf("server listening on :%s", a.cfg.App.Port)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Println("shutdown signal received")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.e.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	} else {
		log.Println("server shutdown gracefully")
	}
}

// Close menutup koneksi DB, Kafka producer, dan Redis.
func (a *App) Close() {
	if a.close != nil {
		a.close()
	}
}
