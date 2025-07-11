package main

import (
	"fmt"
	"log"
	"net/http"

	cfg "tgforwarder/config"
	cfgModel "tgforwarder/config/model"
	httpserver "tgforwarder/http"
	kafka "tgforwarder/messaging"
)

func main() {
	kafkaConf, err := cfg.LoadSection[cfgModel.KafkaConfig]()
	if err != nil {
		log.Fatalf("Kafka config error: %v", err)
	}
	kafka.Init(kafkaConf.BootstrapServersValue)

	serverConf, err := cfg.LoadSection[cfgModel.ServerConfig]()
	if err != nil {
		log.Fatalf("Server config error: %v", err)
	}

	mux := httpserver.SetupRoutes()
	addr := fmt.Sprintf(":%d", serverConf.Port)
	log.Printf("🌐 Webhook server is listening on %s...", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
