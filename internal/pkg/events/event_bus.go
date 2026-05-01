package events

import (
	"context"
	"encoding/json"
	"log/slog"
	"reflect"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	StreamName   = "events_stream"
	GroupName    = "events_group"
	ConsumerName = "event_consumer"
)

// EventBus manages event subscriptions and dispatching using Redis Streams
type EventBus struct {
	redisClient *redis.Client
	handlers    map[string][]EventHandler
	eventTypes  map[string]reflect.Type
	mutex       sync.RWMutex
	logger      *slog.Logger
	ctx         context.Context
	cancel      context.CancelFunc
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
func NewEventBus(redisClient *redis.Client, logger *slog.Logger) *EventBus {
	ctx, cancel := context.WithCancel(context.Background())
	return &EventBus{
		redisClient: redisClient,
		handlers:    make(map[string][]EventHandler),
		eventTypes:  make(map[string]reflect.Type),
		logger:      logger,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Subscribe registers an event handler for a specific event type
func (eb *EventBus) Subscribe(eventType string, handler EventHandler, eventPrototype interface{}) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
	eb.eventTypes[eventType] = reflect.TypeOf(eventPrototype)
}

// Publish sends an event to Redis Stream
func (eb *EventBus) Publish(event Event) {
	eventType := event.GetEventType()

	data, err := json.Marshal(event)
	if err != nil {
		eb.logger.Error("failed to marshal event", slog.String("error", err.Error()))
		return
	}

	err = eb.redisClient.XAdd(eb.ctx, &redis.XAddArgs{
		Stream: StreamName,
		Values: map[string]interface{}{
			"type":    eventType,
			"payload": data,
		},
	}).Err()

	if err != nil {
		eb.logger.Error("failed to publish event to Redis",
			slog.String("type", eventType),
			slog.String("error", err.Error()))
	}
}

// Start initiates the consumer worker
func (eb *EventBus) Start() {
	// Create consumer group if it doesn't exist
	err := eb.redisClient.XGroupCreateMkStream(eb.ctx, StreamName, GroupName, "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		eb.logger.Warn("failed to create consumer group", slog.String("error", err.Error()))
	}

	go eb.worker()
}

// Stop stops the consumer worker
func (eb *EventBus) Stop() {
	eb.cancel()
}

func (eb *EventBus) worker() {
	eb.logger.Info("Event bus worker started")
	for {
		select {
		case <-eb.ctx.Done():
			return
		default:
			streams, err := eb.redisClient.XReadGroup(eb.ctx, &redis.XReadGroupArgs{
				Group:    GroupName,
				Consumer: ConsumerName,
				Streams:  []string{StreamName, ">"},
				Count:    10,
				Block:    time.Second * 2,
			}).Result()

			if err != nil {
				if err != redis.Nil {
					eb.logger.Error("failed to read from redis stream", slog.String("error", err.Error()))
				}
				continue
			}

			for _, stream := range streams {
				for _, msg := range stream.Messages {
					eb.processMessage(msg)
				}
			}
		}
	}
}

func (eb *EventBus) processMessage(msg redis.XMessage) {
	eventType, ok := msg.Values["type"].(string)
	if !ok {
		eb.logger.Warn("message missing event type", slog.String("id", msg.ID))
		return
	}

	payloadStr, ok := msg.Values["payload"].(string)
	if !ok {
		eb.logger.Warn("message missing payload", slog.String("id", msg.ID))
		return
	}

	eb.mutex.RLock()
	handlers, handlersExist := eb.handlers[eventType]
	proto, protoExist := eb.eventTypes[eventType]
	eb.mutex.RUnlock()

	if !handlersExist || !protoExist {
		return
	}

	// Create a new instance of the event type
	eventPtr := reflect.New(proto).Interface()
	if err := json.Unmarshal([]byte(payloadStr), eventPtr); err != nil {
		eb.logger.Error("failed to unmarshal event payload",
			slog.String("type", eventType),
			slog.String("error", err.Error()))
		return
	}

	// The reflect.New returns a pointer, but our handlers might expect the value
	// If the prototype was a struct, reflect.New(proto).Interface() is *Struct
	event := reflect.Indirect(reflect.ValueOf(eventPtr)).Interface().(Event)

	for _, handler := range handlers {
		go func(h EventHandler, e Event) {
			defer func() {
				if r := recover(); r != nil {
					eb.logger.Error("event handler panicked", slog.Any("panic", r))
				}
			}()
			h.Handle(e)
		}(handler, event)
	}

	// Acknowledge the message
	eb.redisClient.XAck(eb.ctx, StreamName, GroupName, msg.ID)
}
