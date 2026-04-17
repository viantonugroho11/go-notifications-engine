package main

import (
	"log"
	"os"

	"github.com/viantonugroho11/go-notifications-engine/internal/bootstrap"
	"github.com/viantonugroho11/go-notifications-engine/internal/config"

	confLoader "github.com/viantonugroho11/go-config-library"
)

func main() {
	cfg := config.Configuration{}
	loader := confLoader.New("", "github.com/viantonugroho11/go-notifications-engine", os.Getenv("CONSUL_URL"),
		confLoader.WithConfigFileSearchPaths("./config"),
	)
	if err := loader.Load(&cfg); err != nil {
		log.Fatalf("config load: %v", err)
	}

	app, err := bootstrap.New(cfg)
	if err != nil {
		log.Fatalf("bootstrap: %v", err)
	}
	defer app.Close()

	app.Run()
}
