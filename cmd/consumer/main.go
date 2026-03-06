package main

import (
	"log"
	"os"

	"go-boilerplate-clean/internal/bootstrap"
	"go-boilerplate-clean/internal/config"

	confLoader "github.com/viantonugroho11/go-config-library"
)

func main() {
	cfg := config.Configuration{}
	loader := confLoader.New("", "go-boilerplate-clean", os.Getenv("CONSUL_URL"),
		confLoader.WithConfigFileSearchPaths("./config"),
	)
	if err := loader.Load(&cfg); err != nil {
		log.Fatalf("config load: %v", err)
	}

	app, err := bootstrap.NewConsumer(cfg)
	if err != nil {
		log.Fatalf("bootstrap consumer: %v", err)
	}
	defer app.Close()

	app.Run()
}
