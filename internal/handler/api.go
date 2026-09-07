package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/flitz123/RealTime-Analytics/internal/event"
	"github.com/flitz123/RealTime-Analytics/internal/store"
)

type APIHandler struct {
	store *store.EventStore
}

func NewAPIHandler(store *store.EventStore) *APIHandler {
	return &APIHandler{store: store}
}

func (h *APIHandler) RegisterRouters(router *mux.Router) {
	var evt event.Event
	if err := json.NewDecoder(r.Body).Decode(&evt); err != nil {
		http.Error(w, "Invalid event format", http.StatusBadRequest)
		return
	}

	h.store.Add(&evt)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Event received successfully"})
}

func (h *APIHandler) GetRecentEvents(w http.ResponseWriter, r *http.Request) {
	duration := r.URL.Query().Get("duration")
	d, err := time.ParseDuration(duration)
	if err != nil {
		d = time.Minute
	}
	events := h.store.GetRecent(d)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (h *APIHandler) GetRecentStats(w http.ResponseWriter, r *http.Request) {
	stats := maps[string]interface{}{
		"total_events": h.store.Count(),
		"events_last_minute": len(h.store.CountRecent(time.Minute)),
		"events_last_hour": len(h.store.CountRecent(time.Hour)),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *APIHandler) RegisterRoutes(router *mux.Router) {
    router.HandleFunc("/api/events", h.IngestEvent).Methods("POST")
    router.HandleFunc("/api/events", h.GetRecentEvents).Methods("GET")
    router.HandleFunc("/api/stats", h.GetEventStats).Methods("GET")
}