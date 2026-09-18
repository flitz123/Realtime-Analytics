package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/flitz123/RealTime-Analytics/internal/analytics"
	"github.com/flitz123/RealTime-Analytics/internal/handler"
	"github.com/flitz123/RealTime-Analytics/internal/store"
)

func main() {
	eventStore := store.NewEventStore()
	processor := analytics.NewProcessor(eventStore)
	processor.Start()

	wsHub := handler.NewWebSocketHub(eventStore, processor)
	apiHandler := handler.NewAPIHandler(eventStore)

	router := mux.NewRouter()

	apiHandler.RegisterRouters(router)

	router.HandleFunc("/ws", wsHub.HandleWebSocket)

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web")))

	log.Println("Server starting on: 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
