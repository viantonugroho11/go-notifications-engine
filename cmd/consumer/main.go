package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"go-boilerplate-clean/internal/bootstrap"
	"go-boilerplate-clean/internal/config"
	brokerkafka "go-boilerplate-clean/internal/infrastructure/broker/kafka"

	confLoader "github.com/viantonugroho11/go-config-library"
)

func main() {
	consumerFlag := flag.String("consumer", "", "satu consumer yang dijalankan. Pilihan: "+strings.Join(brokerkafka.AvailableConsumerKeys(), " | "))
	flag.Parse()

	cfg := config.Configuration{}
	loader := confLoader.New("", "go-boilerplate-clean", os.Getenv("CONSUL_URL"),
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
