package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/viantonugroho11/go-notifications-engine/internal/bootstrap"
	"github.com/viantonugroho11/go-notifications-engine/internal/config"
	brokerkafka "github.com/viantonugroho11/go-notifications-engine/internal/infrastructure/broker/kafka"

	confLoader "github.com/viantonugroho11/go-config-library"
)

func main() {
	consumerFlag := flag.String("consumer", "", "satu consumer yang dijalankan. Pilihan: "+strings.Join(brokerkafka.AvailableConsumerKeys(), " | "))
	flag.Parse()

	cfg := config.Configuration{}
	loader := confLoader.New("", "github.com/viantonugroho11/go-notifications-engine", os.Getenv("CONSUL_URL"),
		confLoader.WithConfigFileSearchPaths("./config"),
	)
	if err := loader.Load(&cfg); err != nil {
		log.Fatalf("config load: %v", err)
	}

	consumerKey := strings.TrimSpace(*consumerFlag)
	app, err := bootstrap.NewConsumer(cfg, consumerKey)
	if err != nil {
		log.Fatalf("bootstrap consumer: %v", err)
	}
	defer app.Close()

	app.Run()
}
