package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/flitz123/Realtime-Analytics/internal/analytics"
	"github.com/flitz123/Realtime-Analytics/internal/handler"
	"github.com/flitz123/Realtime-Analytics/internal/store"
)

func main() {
	eventStore := store.NewEventStore()
	processor := analytics.NewProcessor(eventStore)
	processor.Start()

	wsHub := handler.NewWebSocketHub(eventStore, processor)
	apiHandler := handler.NewAPIHandler(eventStore)

	router := mux.NewRouter()

	apiHandler.RegisterRoutes(router)

	router.HandleFunc("/ws", wsHub.HandleWebSocket)

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server starting on port:", port)
	if err := http.ListenAndServe(":"+port, withCORS(router)); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := os.Getenv("FRONTEND_URL")
		if origin == "" {
			origin = "*"
		}
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
