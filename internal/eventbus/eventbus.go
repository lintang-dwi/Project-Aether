package eventbus

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

// Subscription represents a single subscriber to a topic.
type Subscription struct {
	ID       string
	Topic    string
	Handler  EventHandler
	ch       chan Event
	cancel   context.CancelFunc
}

// EventBus is an in-process pub/sub event bus supporting topic pattern matching.
type EventBus struct {
	mu            sync.RWMutex
	subscribers   map[string][]*Subscription
	bufferSize    int
	isClosed      bool
}

// NewEventBus creates a new EventBus with the specified channel buffer size.
func NewEventBus(bufferSize int) *EventBus {
	if bufferSize <= 0 {
		bufferSize = 256
	}
	return &EventBus{
		subscribers: make(map[string][]*Subscription),
		bufferSize:  bufferSize,
	}
}

// Subscribe registers a handler for events matching the specified topic pattern.
// Pattern can be exact (e.g. "workspace.changed") or prefix wildcard (e.g. "workspace.*" or "*").
func (eb *EventBus) Subscribe(topicPattern string, handler EventHandler) *Subscription {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if eb.isClosed {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	sub := &Subscription{
		ID:      fmt.Sprintf("%s-%d", topicPattern, len(eb.subscribers[topicPattern])),
		Topic:   topicPattern,
		Handler: handler,
		ch:      make(chan Event, eb.bufferSize),
		cancel:  cancel,
	}

	eb.subscribers[topicPattern] = append(eb.subscribers[topicPattern], sub)

	// Start worker loop for this subscriber
	go func(s *Subscription, ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-s.ch:
				if !ok {
					return
				}
				s.Handler(event)
			}
		}
	}(sub, ctx)

	return sub
}

// Publish broadcasts an event asynchronously to all subscribers matching the event's topic.
func (eb *EventBus) Publish(event Event) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	if eb.isClosed {
		return
	}

	for pattern, subs := range eb.subscribers {
		if matchTopic(pattern, event.Topic) {
			for _, sub := range subs {
				select {
				case sub.ch <- event:
				default:
					// Non-blocking publish if buffer is full to prevent deadlocks
				}
			}
		}
	}
}

// Unsubscribe removes a subscription and closes its channel.
func (eb *EventBus) Unsubscribe(sub *Subscription) {
	if sub == nil {
		return
	}

	eb.mu.Lock()
	defer eb.mu.Unlock()

	subs, exists := eb.subscribers[sub.Topic]
	if !exists {
		return
	}

	for i, s := range subs {
		if s == sub {
			sub.cancel()
			close(sub.ch)
			eb.subscribers[sub.Topic] = append(subs[:i], subs[i+1:]...)
			break
		}
	}
}

// Close gracefully closes the event bus and all active subscriptions.
func (eb *EventBus) Close() {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if eb.isClosed {
		return
	}

	eb.isClosed = true
	for _, subs := range eb.subscribers {
		for _, sub := range subs {
			sub.cancel()
			close(sub.ch)
		}
	}
	eb.subscribers = make(map[string][]*Subscription)
}

// matchTopic checks if a topic matches a pattern.
// Supports "*", "workspace.*", and exact match.
func matchTopic(pattern, topic string) bool {
	if pattern == "*" || pattern == topic {
		return true
	}
	if strings.HasSuffix(pattern, ".*") {
		prefix := strings.TrimSuffix(pattern, ".*")
		return strings.HasPrefix(topic, prefix+".") || topic == prefix
	}
	return false
}
