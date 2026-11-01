package broker

import (
	"context"
	"log"
	"sync"
)

type subscription[T any] struct {
	topics []string
	ch     chan T
}

type brokerType[T any] struct {
	topic   string
	content T
}

type Broker[T any] struct {
	register   chan subscription[T]
	unregister chan subscription[T]
	publish    chan brokerType[T]
	clients    map[string]map[chan T]struct{}
	subs       map[chan T][]string
	done       chan struct{}
}

func NewBroker[T any]() *Broker[T] {
	return &Broker[T]{
		register:   make(chan subscription[T]),
		unregister: make(chan subscription[T]),
		publish:    make(chan brokerType[T], 64),
		clients:    make(map[string]map[chan T]struct{}),
		subs:       make(map[chan T][]string),
		done:       make(chan struct{}),
	}
}

// pintu depan buat controller
func (b *Broker[T]) Subscribe(topics []string) (<-chan T, func()) {
	ch := make(chan T, 8)

	seen := make(map[string]struct{}, len(topics))
	ids := make([]string, 0, len(topics))
	for _, topic := range topics {
		if _, ok := seen[topic]; ok {
			continue
		}
		seen[topic] = struct{}{}
		ids = append(ids, topic)
	}

	select {
	case b.register <- subscription[T]{topics: ids, ch: ch}:
	case <-b.done:
		// broker sudah mati: kasih channel yang langsung tertutup
		close(ch)
		return ch, func() {}
	}

	var once sync.Once
	unsubscribe := func() {
		once.Do(func() {
			select {
			case b.unregister <- subscription[T]{ch: ch}:
			case <-b.done:
			}
		})
	}
	return ch, unsubscribe
}

// pintu depan buat service
func (b *Broker[T]) Publish(topic string, content T) {
	select {
	case b.publish <- brokerType[T]{topic: topic, content: content}:
	default:
		log.Printf("broker: publish buffer penuh, update topic %s dibuang", topic)
	}
}

func (b *Broker[T]) add(s subscription[T]) {
	for _, topic := range s.topics {
		if b.clients[topic] == nil {
			b.clients[topic] = make(map[chan T]struct{})
		}
		b.clients[topic][s.ch] = struct{}{}
	}
	b.subs[s.ch] = s.topics
}

func (b *Broker[T]) remove(ch chan T) {
	ids, ok := b.subs[ch]
	if !ok {
		return // belum pernah terdaftar / sudah dihapus
	}
	for _, topic := range ids {
		delete(b.clients[topic], ch)
		if len(b.clients[topic]) == 0 {
			delete(b.clients, topic)
		}
	}
	delete(b.subs, ch)
	close(ch)
}

func (b *Broker[T]) Run(ctx context.Context) {
	defer func() {
		close(b.done)
		for ch := range b.subs {
			close(ch)
		}
	}()

	for {
		select {
		case sub := <-b.register:
			b.add(sub)
		case sub := <-b.unregister:
			b.remove(sub.ch)
		case pub := <-b.publish:
			for ch := range b.clients[pub.topic] {
				select {
				case ch <- pub.content:
				default:
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
