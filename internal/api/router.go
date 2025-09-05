package httpserver

import (
	"net/http"

	"io"
	"log"
	kafka "tgforwarder/internal/messaging"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/webhook", Webhook)
	return mux
}

func Webhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = kafka.Send("telegram-updates", body)
	if err != nil {
		log.Printf("Error delivering message to Kafka: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("Message delivered to Kafka")
	w.WriteHeader(http.StatusOK)
}
