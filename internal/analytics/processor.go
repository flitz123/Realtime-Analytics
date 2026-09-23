package analytics

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/flitz123/Realtime-Analytics/internal/events"
	"github.com/flitz123/Realtime-Analytics/internal/store"
)

type Metrics struct {
	TotalEvents           int                      `json:"total_events"`
	EventsPerMinute       float64                  `json:"events_per_minute"`
	EventTypeDistribution map[events.EventType]int `json:"event_type_distribution"`
	ActiveUsers           int                      `json:"active_users"`
	Timestamp             time.Time                `json:"timestamp"`
}

type Processor struct {
	store       *store.EventStore
	metrics     *Metrics
	mu          sync.RWMutex
	subscribers map[string]chan []byte
}

func NewProcessor(store *store.EventStore) *Processor {
	return &Processor{
		store:       store,
		metrics:     &Metrics{},
		subscribers: make(map[string]chan []byte),
	}
}

func (p *Processor) Start() {
	ticker := time.NewTicker(5 * time.Second)

	go func() {
		for range ticker.C {
			p.updateMetrics()
		}
	}()
}

func (p *Processor) updateMetrics() {
	p.mu.Lock()
	defer p.mu.Unlock()

	recentEvents := p.store.GetRecent(time.Minute)
	p.metrics.TotalEvents = p.store.GetCount()
	p.metrics.EventsPerMinute = float64(len(recentEvents))
	p.metrics.Timestamp = time.Now()
	distribution := make(map[events.EventType]int)
	uniqueUsers := make(map[string]bool)

	for _, e := range recentEvents {
		distribution[e.Type]++
		uniqueUsers[e.UserID] = true
	}

	p.metrics.EventTypeDistribution = distribution
	p.metrics.ActiveUsers = len(uniqueUsers)
	p.broadcastMetrics()
}

func (p *Processor) GetMetrics() *Metrics {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.metrics
}

func (p *Processor) Subscribe(id string) chan []byte {
	p.mu.Lock()
	defer p.mu.Unlock()

	ch := make(chan []byte, 10)
	p.subscribers[id] = ch
	return ch
}

func (p *Processor) Unsubscribe(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if ch, ok := p.subscribers[id]; ok {
		close(ch)
		delete(p.subscribers, id)
	}
}

func (p *Processor) broadcastMetrics() {
	metricsJSON, err := json.Marshal(p.metrics)
	if err != nil {
		return
	}

	for _, ch := range p.subscribers {
		select {
		case ch <- metricsJSON:
		default:
		}
	}
}
