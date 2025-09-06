package api

import (
	"io"
	"net/http"
	"tgforwarder/internal/logger"
	"tgforwarder/internal/messaging"
)

type Options struct {
	Logger          logger.Logger
	MessProducer    *messaging.KafkaProducer
	TgMessTopicName string
	TgWebHookPath   string
	HealthCheckPath string
}

func SetupRoutes(opt Options) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc(opt.TgWebHookPath, func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			opt.Logger.Error("Failed to read request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if err := opt.MessProducer.Send(r.Context(), opt.TgMessTopicName, body); err != nil {
			opt.Logger.Error("Error delivering message to Kafka: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		opt.Logger.Info("Message delivered to Kafka")
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc(opt.HealthCheckPath, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	return mux
}
