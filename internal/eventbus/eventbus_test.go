package eventbus_test

import (
	"sync"
	"testing"
	"time"

	"aether/internal/eventbus"
)

func TestEventBus_PubSubExact(t *testing.T) {
	eb := eventbus.NewEventBus(100)
	defer eb.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	var received eventbus.Event
	eb.Subscribe("workspace.changed", func(e eventbus.Event) {
		received = e
		wg.Done()
	})

	evt := eventbus.NewEvent("workspace.changed", "test", map[string]interface{}{"path": "main.go"})
	eb.Publish(evt)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		if received.Topic != "workspace.changed" {
			t.Errorf("expected topic 'workspace.changed', got '%s'", received.Topic)
		}
		if received.Payload["path"] != "main.go" {
			t.Errorf("expected path 'main.go', got '%v'", received.Payload["path"])
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for event")
	}
}

func TestEventBus_WildcardMatching(t *testing.T) {
	eb := eventbus.NewEventBus(100)
	defer eb.Close()

	var count int32
	var mu sync.Mutex

	eb.Subscribe("workspace.*", func(e eventbus.Event) {
		mu.Lock()
		count++
		mu.Unlock()
	})

	eb.Publish(eventbus.NewEvent("workspace.file_added", "test", nil))
	eb.Publish(eventbus.NewEvent("workspace.file_deleted", "test", nil))
	eb.Publish(eventbus.NewEvent("other.topic", "test", nil))

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if count != 2 {
		t.Errorf("expected 2 wildcard events, got %d", count)
	}
}

func TestEventBus_Unsubscribe(t *testing.T) {
	eb := eventbus.NewEventBus(100)
	defer eb.Close()

	var count int
	var mu sync.Mutex

	sub := eb.Subscribe("test.topic", func(e eventbus.Event) {
		mu.Lock()
		count++
		mu.Unlock()
	})

	eb.Publish(eventbus.NewEvent("test.topic", "test", nil))
	time.Sleep(50 * time.Millisecond)

	eb.Unsubscribe(sub)
	eb.Publish(eventbus.NewEvent("test.topic", "test", nil))
	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if count != 1 {
		t.Errorf("expected 1 event after unsubscribe, got %d", count)
	}
}
