package store

import (
	"sync"
	"time"

	"github.com/flitz123/RealTime-Analytics/internal/event"
)

type EventStore struct {
	mu     sync.RWMutex
	events []event.Event
}

func NewEventStore() *EventStore {
	return &EventStore{
		events: make([]event.Event, 0),
	}
}

func (s *EventStore) Add(event event.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = append(s.events, event)

	if len(s.events) > 1000 {
		s.events = s.events[len(s.events)-1000:]
	}
}

func (s *EventStore) GetRecent(duration time.Duration) []event.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cutoff := time.Now().Add(-duration)
	var recent []event.Event

	for i := len(s.events) - 1; i >= 0; i-- {
		if s.events[i].Timestamp.After(cutoff) {
			recent = append(recent, s.events[i])
		}
	}
	return recent
}

func (s *EventStore) GetByType(eventType event.EventType, duration time.Duration) []event.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cutoff := time.Now().Add(-duration)
	var filtered []event.Event

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
