package events

import (
	"sync"
	"testing"
)

func TestNewBus(t *testing.T) {
	b := NewBus()
	if b == nil {
		t.Fatal("expected non-nil bus")
	}
	if len(b.Topics()) != 0 {
		t.Fatalf("expected 0 topics, got %d", len(b.Topics()))
	}
}

func TestPublishSubscribe(t *testing.T) {
	b := NewBus()
	received := false
	var receivedPayload any

	err := b.Subscribe("test", func(e Event) {
		received = true
		receivedPayload = e.Payload
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = b.Publish("test", "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !received {
		t.Fatal("expected handler to be called")
	}
	if receivedPayload != "hello" {
		t.Fatalf("expected payload 'hello', got %v", receivedPayload)
	}
}

func TestSubscribeNilHandler(t *testing.T) {
	b := NewBus()
	err := b.Subscribe("test", nil)
	if err == nil {
		t.Fatal("expected error for nil handler")
	}
}

func TestPublishNoSubscribers(t *testing.T) {
	b := NewBus()
	err := b.Publish("no-subscribers", "data")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMultipleSubscribers(t *testing.T) {
	b := NewBus()
	count := 0
	var mu sync.Mutex

	for i := 0; i < 3; i++ {
		_ = b.Subscribe("test", func(e Event) {
			mu.Lock()
			count++
			mu.Unlock()
		})
	}

	_ = b.Publish("test", nil)

	mu.Lock()
	defer mu.Unlock()
	if count != 3 {
		t.Fatalf("expected 3 handler calls, got %d", count)
	}
}

func TestUnsubscribe(t *testing.T) {
	b := NewBus()
	_ = b.Subscribe("test", func(e Event) {})

	err := b.Unsubscribe("test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.SubscriberCount("test") != 0 {
		t.Fatalf("expected 0 subscribers after unsubscribe, got %d", b.SubscriberCount("test"))
	}
}

func TestUnsubscribeNotFound(t *testing.T) {
	b := NewBus()
	err := b.Unsubscribe("missing")
	if err == nil {
		t.Fatal("expected error for unsubscribing from non-existent topic")
	}
}

func TestTopics(t *testing.T) {
	b := NewBus()
	_ = b.Subscribe("a", func(e Event) {})
	_ = b.Subscribe("b", func(e Event) {})
	_ = b.Subscribe("c", func(e Event) {})

	topics := b.Topics()
	if len(topics) != 3 {
		t.Fatalf("expected 3 topics, got %d", len(topics))
	}
}

func TestSubscriberCount(t *testing.T) {
	b := NewBus()
	_ = b.Subscribe("test", func(e Event) {})
	_ = b.Subscribe("test", func(e Event) {})

	if b.SubscriberCount("test") != 2 {
		t.Fatalf("expected 2 subscribers, got %d", b.SubscriberCount("test"))
	}
	if b.SubscriberCount("other") != 0 {
		t.Fatalf("expected 0 subscribers for other topic, got %d", b.SubscriberCount("other"))
	}
}

func TestPublishMultipleTopics(t *testing.T) {
	b := NewBus()
	var received []string
	var mu sync.Mutex

	_ = b.Subscribe("topic-a", func(e Event) {
		mu.Lock()
		received = append(received, "a:"+e.Payload.(string))
		mu.Unlock()
	})
	_ = b.Subscribe("topic-b", func(e Event) {
		mu.Lock()
		received = append(received, "b:"+e.Payload.(string))
		mu.Unlock()
	})

	_ = b.Publish("topic-a", "1")
	_ = b.Publish("topic-b", "2")

	mu.Lock()
	defer mu.Unlock()
	if len(received) != 2 {
		t.Fatalf("expected 2 received events, got %d", len(received))
	}
}
