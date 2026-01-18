package events

import (
	"sync"
)

// EventBus manages event subscriptions and dispatching
type EventBus struct {
	handlers map[string][]EventHandler
	mutex    sync.RWMutex
}

// Event interface defines the contract for events
type Event interface {
	GetEventType() string
}

// EventHandler interface defines the contract for event handlers
type EventHandler interface {
	Handle(event Event)
}

// NewEventBus creates a new event bus instance
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]EventHandler),
	}
}

// Subscribe registers an event handler for a specific event type
func (eb *EventBus) Subscribe(eventType string, handler EventHandler) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
}

// Publish sends an event to all registered handlers for that event type
func (eb *EventBus) Publish(event Event) {
	eventType := event.GetEventType()

	eb.mutex.RLock()
	handlers, exists := eb.handlers[eventType]
	eb.mutex.RUnlock()

	if !exists {
		return
	}

	// Dispatch to all handlers concurrently
	var wg sync.WaitGroup
	for _, handler := range handlers {
		wg.Add(1)
		go func(h EventHandler) {
			defer wg.Done()
			h.Handle(event)
		}(handler)
	}
	wg.Wait()
}
