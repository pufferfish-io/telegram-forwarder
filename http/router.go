package httpserver

import (
	"net/http"
	"tgforwarder/http/handler"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/webhook", handler.Webhook)
	return mux
}
