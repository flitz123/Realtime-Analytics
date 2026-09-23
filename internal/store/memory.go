package store

import (
	"sync"
	"time"

	"github.com/flitz123/Realtime-Analytics/internal/events"
)

type EventStore struct {
	mu     sync.RWMutex
	events []events.Event
}

func NewEventStore() *EventStore {
	return &EventStore{
		events: make([]events.Event, 0),
	}
}

func (s *EventStore) Add(event events.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = append(s.events, event)

	if len(s.events) > 1000 {
		s.events = s.events[len(s.events)-1000:]
	}
}

func (s *EventStore) GetRecent(duration time.Duration) []events.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cutoff := time.Now().Add(-duration)
	var recent []events.Event

	for i := len(s.events) - 1; i >= 0; i-- {
		if s.events[i].Timestamp.After(cutoff) {
			recent = append(recent, s.events[i])
		}
	}
	return recent
}

func (s *EventStore) GetByType(eventType events.EventType, duration time.Duration) []events.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cutoff := time.Now().Add(-duration)
	var filtered []events.Event

	for i := len(s.events) - 1; i >= 0; i-- {
		if s.events[i].Type == eventType && s.events[i].Timestamp.After(cutoff) {
			filtered = append(filtered, s.events[i])
		}
	}
	return filtered
}

func (s *EventStore) GetCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.events)
}
