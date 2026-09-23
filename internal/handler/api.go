package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/flitz123/Realtime-Analytics/internal/events"
	"github.com/flitz123/Realtime-Analytics/internal/store"
)

type APIHandler struct {
	store *store.EventStore
}

func NewAPIHandler(store *store.EventStore) *APIHandler {
	return &APIHandler{store: store}
}

func (h *APIHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/events", h.IngestEvent).Methods(http.MethodPost)
	router.HandleFunc("/api/events", h.GetRecentEvents).Methods(http.MethodGet)
	router.HandleFunc("/api/stats", h.GetRecentStats).Methods(http.MethodGet)
}

func (h *APIHandler) IngestEvent(w http.ResponseWriter, r *http.Request) {
	var evt events.Event
	if err := json.NewDecoder(r.Body).Decode(&evt); err != nil {
		http.Error(w, "Invalid event format", http.StatusBadRequest)
		return
	}

	if evt.Timestamp.IsZero() {
		evt.Timestamp = time.Now()
	}
	h.store.Add(evt)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Event received successfully"})
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
	stats := map[string]interface{}{
		"total_events":       h.store.GetCount(),
		"events_last_minute": len(h.store.GetRecent(time.Minute)),
		"events_last_hour":   len(h.store.GetRecent(time.Hour)),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
