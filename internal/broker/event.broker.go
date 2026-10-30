package broker

import (
	"context"
	"log"
	"rcontrisha/backend-eventhub/internal/dto"
	"sync"
)

type subscription struct {
	eventIDs []string
	ch       chan dto.UpdateAttendee
}

type Broker struct {
	register   chan subscription
	unregister chan subscription
	publish    chan dto.UpdateAttendee
	clients    map[string]map[chan dto.UpdateAttendee]struct{}
	subs       map[chan dto.UpdateAttendee][]string
	done       chan struct{}
}

func NewBroker() *Broker {
	return &Broker{
		register:   make(chan subscription),
		unregister: make(chan subscription),
		publish:    make(chan dto.UpdateAttendee, 64),
		clients:    make(map[string]map[chan dto.UpdateAttendee]struct{}),
		subs:       make(map[chan dto.UpdateAttendee][]string),
		done:       make(chan struct{}),
	}
}

// pintu depan buat controller
func (b *Broker) Subscribe(eventIDs []string) (<-chan dto.UpdateAttendee, func()) {
	ch := make(chan dto.UpdateAttendee, 8)

	// salin + dedup, supaya slice milik pemanggil nggak kebawa ke dalam broker
	seen := make(map[string]struct{}, len(eventIDs))
	ids := make([]string, 0, len(eventIDs))
	for _, id := range eventIDs {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}

	select {
	case b.register <- subscription{eventIDs: ids, ch: ch}:
	case <-b.done:
		// broker sudah mati: kasih channel yang langsung tertutup
		close(ch)
		return ch, func() {}
	}

	var once sync.Once
	unsubscribe := func() {
		once.Do(func() {
			select {
			case b.unregister <- subscription{ch: ch}:
			case <-b.done:
			}
		})
	}
	return ch, unsubscribe
}

// pintu depan buat service
func (b *Broker) Publish(u dto.UpdateAttendee) {
	// TODO: kirim ke b.publish
	select {
	case b.publish <- u:
	default:
		log.Printf("broker: publish buffer penuh, update event %s dibuang", u.EventID)
	}
}

func (b *Broker) add(s subscription) {
	for _, id := range s.eventIDs {
		if b.clients[id] == nil {
			b.clients[id] = make(map[chan dto.UpdateAttendee]struct{})
		}
		b.clients[id][s.ch] = struct{}{}
	}
	b.subs[s.ch] = s.eventIDs
}

func (b *Broker) remove(ch chan dto.UpdateAttendee) {
	ids, ok := b.subs[ch]
	if !ok {
		return // belum pernah terdaftar / sudah dihapus
	}
	for _, id := range ids {
		delete(b.clients[id], ch)
		if len(b.clients[id]) == 0 {
			delete(b.clients, id)
		}
	}
	delete(b.subs, ch)
	close(ch)
}

func (b *Broker) Run(ctx context.Context) {
	defer func() {
		close(b.done)
		for ch := range b.subs {
			close(ch)
		}
	}()

	for {
		select {
		case s := <-b.register:
			b.add(s)
		case s := <-b.unregister:
			b.remove(s.ch)
		case u := <-b.publish:
			for ch := range b.clients[u.EventID] {
				select {
				case ch <- u:
				default:
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
