package main

import (
	"log"
	"net/http"

	httpserver "tgforwarder/internal/api"
	"tgforwarder/internal/config"
	kafka "tgforwarder/internal/messaging"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	kafka.Init(kafka.Option{BootstrapServersKey: cfg.Kafka.BootstrapServersKey, Brokers: cfg.Kafka.BootstrapServersValue})

	mux := httpserver.SetupRoutes()
	log.Printf("🌐 Webhook server is listening on %s...", cfg.Server.Addr)
	log.Fatal(http.ListenAndServe(cfg.Server.Addr, mux))
}
