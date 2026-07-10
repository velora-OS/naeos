package events

import (
	"fmt"
	"sync"
)

type Event struct {
	Topic   string
	Payload any
}

type Handler func(event Event)

type EventBus interface {
	Publish(topic string, payload any) error
	Subscribe(topic string, handler Handler) error
	Unsubscribe(topic string) error
	Topics() []string
	SubscriberCount(topic string) int
}

type Bus struct {
	mu          sync.RWMutex
	subscribers map[string][]Handler
}

func NewBus() *Bus {
	return &Bus{
		subscribers: make(map[string][]Handler),
	}
}

func (b *Bus) Publish(topic string, payload any) error {
	b.mu.RLock()
	handlers := b.subscribers[topic]
	b.mu.RUnlock()

	event := Event{Topic: topic, Payload: payload}
	for _, handler := range handlers {
		handler(event)
	}
	return nil
}

func (b *Bus) Subscribe(topic string, handler Handler) error {
	if handler == nil {
		return fmt.Errorf("handler must not be nil")
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers[topic] = append(b.subscribers[topic], handler)
	return nil
}

func (b *Bus) Unsubscribe(topic string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, exists := b.subscribers[topic]; !exists {
		return fmt.Errorf("no subscribers for topic %s", topic)
	}
	delete(b.subscribers, topic)
	return nil
}

func (b *Bus) Topics() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	result := make([]string, 0, len(b.subscribers))
	for topic := range b.subscribers {
		result = append(result, topic)
	}
	return result
}

func (b *Bus) SubscriberCount(topic string) int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.subscribers[topic])
}
