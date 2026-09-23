package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/flitz123/Realtime-Analytics/internal/analytics"
	"github.com/flitz123/Realtime-Analytics/internal/events"
	"github.com/flitz123/Realtime-Analytics/internal/store"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketHub struct {
	clients   map[*websocket.Conn]bool
	mu        sync.Mutex
	store     *store.EventStore
	processor *analytics.Processor
}

func NewWebSocketHub(store *store.EventStore, processor *analytics.Processor) *WebSocketHub {
	return &WebSocketHub{
		store:     store,
		processor: processor,
		clients:   make(map[*websocket.Conn]bool),
	}
}

func (h *WebSocketHub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	h.mu.Lock()
	h.clients[conn] = true
	h.mu.Unlock()

	subscriberID := uuid.New().String()
	metricsChan := h.processor.Subscribe(subscriberID)

	defer func() {
		h.mu.Lock()
		delete(h.clients, conn)
		h.mu.Unlock()
		h.processor.Unsubscribe(subscriberID)
		conn.Close()
	}()

	go func() {
		for metrics := range metricsChan {
			if err := conn.WriteMessage(websocket.TextMessage, metrics); err != nil {
				log.Printf("Failed to send metrics: %v", err)
				return
			}
		}
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var evt events.Event
		if err := json.Unmarshal(message, &evt); err != nil {
			log.Printf("Failed to unmarshal event: %v", err)
			continue
		}

		h.store.Add(evt)

		h.broadcastEvent(&evt)
	}
}

func (h *WebSocketHub) broadcastEvent(evt *events.Event) {
	eventJSON, err := evt.ToJSON()
	if err != nil {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for client := range h.clients {
		if err := client.WriteMessage(websocket.TextMessage, eventJSON); err != nil {
			log.Printf("Failed to send event to client: %v", err)
			client.Close()
			delete(h.clients, client)
		}
	}
}
