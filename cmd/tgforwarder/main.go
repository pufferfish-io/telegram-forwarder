package main

import (
	"log"
	"net/http"

	"tgforwarder/internal/api"
	"tgforwarder/internal/config"
	"tgforwarder/internal/logger"
	"tgforwarder/internal/messaging"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	logger, clean := logger.NewZapLogger()
	defer clean()

    prod, err := messaging.NewKafkaProducer(messaging.Option{
        Logger:       logger,
        Broker:       cfg.Kafka.BootstrapServersValue,
        SaslUsername: cfg.Kafka.SaslUsername,
        SaslPassword: cfg.Kafka.SaslPassword,
    })

	if err != nil {
		log.Fatalf("config: %v", err)
	}
	defer prod.Close()

	mux := api.SetupRoutes(api.Options{Logger: logger, MessProducer: prod, TgMessTopicName: cfg.Kafka.TgMessTopicName, TgWebHookPath: cfg.Api.TgWebHookPath, HealthCheckPath: cfg.Api.HealthCheckPath})
	log.Printf("🌐 Webhook server is listening on %s...", cfg.Server.Addr)
	log.Fatal(http.ListenAndServe(cfg.Server.Addr, mux))
}
