package messagequeue

import (
	"sync"
	"testing"
	"time"
)

func TestNewQueue(t *testing.T) {
	q := NewQueue("test", 100)
	if q == nil {
		t.Fatal("expected queue to be created")
	}
	if q.Name() != "test" {
		t.Errorf("expected name 'test', got %s", q.Name())
	}
}

func TestPublish(t *testing.T) {
	q := NewQueue("test", 100)

	msg := &Message{
		ID:      "1",
		Topic:   "test",
		Payload: "hello",
	}

	err := q.Publish(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if q.Len() != 1 {
		t.Errorf("expected length 1, got %d", q.Len())
	}
}

func TestPublishFull(t *testing.T) {
	q := NewQueue("test", 1)

	msg1 := &Message{ID: "1", Payload: "a"}
	msg2 := &Message{ID: "2", Payload: "b"}

	q.Publish(msg1)
	err := q.Publish(msg2)
	if err != ErrQueueFull {
		t.Errorf("expected ErrQueueFull, got %v", err)
	}
}

func TestSubscribe(t *testing.T) {
	q := NewQueue("test", 100)

	var received *Message
	var wg sync.WaitGroup
	wg.Add(1)

	q.Subscribe(func(msg *Message) error {
		received = msg
		wg.Done()
		return nil
	})

	msg := &Message{ID: "1", Payload: "hello"}
	q.Publish(msg)

	wg.Wait()

	if received == nil {
		t.Fatal("expected message to be received")
	}
	if received.Payload != "hello" {
		t.Errorf("expected 'hello', got %v", received.Payload)
	}
}

func TestSubscribeRetry(t *testing.T) {
	q := NewQueue("test", 100)

	var attempts int
	var wg sync.WaitGroup
	wg.Add(1)

	q.Subscribe(func(msg *Message) error {
		attempts++
		if attempts < 3 {
			return &QueueError{"retry"}
		}
		wg.Done()
		return nil
	})

	msg := &Message{ID: "1", Payload: "hello", MaxRetries: 3}
	q.Publish(msg)

	wg.Wait()

	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestStop(t *testing.T) {
	q := NewQueue("test", 100)
	q.Subscribe(func(msg *Message) error {
		return nil
	})

	q.Stop()

	q.mu.RLock()
	running := q.running
	q.mu.RUnlock()

	if running {
		t.Error("expected queue to be stopped")
	}
}

func TestTopic(t *testing.T) {
	topic := NewTopic("events")
	if topic == nil {
		t.Fatal("expected topic to be created")
	}

	var received *Message
	var wg sync.WaitGroup
	wg.Add(1)

	topic.Subscribe("sub1", func(msg *Message) error {
		received = msg
		wg.Done()
		return nil
	})

	msg := &Message{ID: "1", Payload: "event"}
	topic.Publish(msg)

	wg.Wait()

	if received == nil {
		t.Fatal("expected message to be received")
	}
}

func TestTopicMultipleSubscribers(t *testing.T) {
	topic := NewTopic("events")

	var count int
	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(2)

	topic.Subscribe("sub1", func(msg *Message) error {
		mu.Lock()
		count++
		mu.Unlock()
		wg.Done()
		return nil
	})

	topic.Subscribe("sub2", func(msg *Message) error {
		mu.Lock()
		count++
		mu.Unlock()
		wg.Done()
		return nil
	})

	msg := &Message{ID: "1", Payload: "event"}
	topic.Publish(msg)

	wg.Wait()

	if count != 2 {
		t.Errorf("expected 2 receives, got %d", count)
	}
}

func TestTopicUnsubscribe(t *testing.T) {
	topic := NewTopic("events")

	topic.Subscribe("sub1", func(msg *Message) error {
		return nil
	})

	if topic.Subscribers() != 1 {
		t.Errorf("expected 1 subscriber, got %d", topic.Subscribers())
	}

	topic.Unsubscribe("sub1")

	if topic.Subscribers() != 0 {
		t.Errorf("expected 0 subscribers, got %d", topic.Subscribers())
	}
}

func TestBroker(t *testing.T) {
	broker := NewBroker()

	broker.CreateTopic("events")
	broker.CreateTopic("logs")

	topics := broker.ListTopics()
	if len(topics) != 2 {
		t.Errorf("expected 2 topics, got %d", len(topics))
	}
}

func TestBrokerPublish(t *testing.T) {
	broker := NewBroker()

	var received *Message
	var wg sync.WaitGroup
	wg.Add(1)

	broker.Subscribe("events", "sub1", func(msg *Message) error {
		received = msg
		wg.Done()
		return nil
	})

	msg := &Message{ID: "1", Payload: "hello"}
	broker.Publish("events", msg)

	wg.Wait()

	if received == nil {
		t.Fatal("expected message to be received")
	}
}

func TestBrokerPublishNotFound(t *testing.T) {
	broker := NewBroker()

	msg := &Message{ID: "1", Payload: "hello"}
	err := broker.Publish("nonexistent", msg)
	if err != ErrTopicNotFound {
		t.Errorf("expected ErrTopicNotFound, got %v", err)
	}
}

func TestBrokerDeleteTopic(t *testing.T) {
	broker := NewBroker()

	broker.CreateTopic("events")
	broker.DeleteTopic("events")

	topics := broker.ListTopics()
	if len(topics) != 0 {
		t.Errorf("expected 0 topics, got %d", len(topics))
	}
}

func TestBrokerStop(t *testing.T) {
	broker := NewBroker()

	broker.CreateTopic("events")
	broker.CreateTopic("logs")

	broker.Stop()

	topics := broker.ListTopics()
	if len(topics) != 2 {
		t.Errorf("expected 2 topics still listed, got %d", len(topics))
	}
}

func TestNewMessage(t *testing.T) {
	msg := NewMessage("topic", "payload")
	if msg == nil {
		t.Fatal("expected message to be created")
	}
	if msg.Topic != "topic" {
		t.Errorf("expected topic 'topic', got %s", msg.Topic)
	}
	if msg.Payload != "payload" {
		t.Errorf("expected payload 'payload', got %v", msg.Payload)
	}
	if msg.MaxRetries != 3 {
		t.Errorf("expected max retries 3, got %d", msg.MaxRetries)
	}
}

func TestMessageTimestamp(t *testing.T) {
	before := time.Now()
	msg := NewMessage("topic", "payload")
	after := time.Now()

	if msg.Timestamp.Before(before) || msg.Timestamp.After(after) {
		t.Error("expected timestamp between before and after")
	}
}
