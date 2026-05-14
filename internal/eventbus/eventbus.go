// Package eventbus provides a lightweight in-process publish/subscribe
// mechanism for broadcasting drift events between driftwatch components.
//
// Subscribers register interest in named topics and receive copies of
// published payloads on a buffered channel. Slow consumers are not blocked;
// if a subscriber's buffer is full the delivery is silently dropped and the
// drop is counted so callers can observe back-pressure.
package eventbus

import (
	"sync"
)

// Event carries a topic name and an opaque payload.
type Event struct {
	Topic   string
	Payload any
}

// Subscription is returned by Subscribe. Close it to stop receiving events.
type Subscription struct {
	C      <-chan Event
	ch     chan Event
	topic  string
	bus    *Bus
	drops  uint64
	mu     sync.Mutex
}

// Drops returns the number of events that were discarded because the
// subscriber's buffer was full at the time of publication.
func (s *Subscription) Drops() uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.drops
}

// Close unregisters the subscription and drains its channel.
func (s *Subscription) Close() {
	s.bus.unsubscribe(s)
}

// Bus is a simple topic-based event broker. The zero value is not usable;
// create one with New.
type Bus struct {
	mu   sync.RWMutex
	subs map[string][]*Subscription
}

// New returns an initialised Bus.
func New() *Bus {
	return &Bus{
		subs: make(map[string][]*Subscription),
	}
}

// Subscribe registers interest in the given topic. bufSize controls how many
// events may be queued before drops occur. A bufSize of zero is allowed but
// will cause drops unless the caller drains the channel immediately.
func (b *Bus) Subscribe(topic string, bufSize int) *Subscription {
	ch := make(chan Event, bufSize)
	s := &Subscription{
		C:     ch,
		ch:    ch,
		topic: topic,
		bus:   b,
	}

	b.mu.Lock()
	b.subs[topic] = append(b.subs[topic], s)
	b.mu.Unlock()

	return s
}

// Publish sends ev to every subscriber registered for ev.Topic.
// It never blocks; subscribers that cannot keep up will have events dropped.
func (b *Bus) Publish(ev Event) {
	b.mu.RLock()
	list := b.subs[ev.Topic]
	// Copy the slice so we can release the read lock before sending.
	copy := make([]*Subscription, len(list))
	copy_n := copy_slice(copy, list)
	b.mu.RUnlock()

	for _, s := range copy[:copy_n] {
		select {
		case s.ch <- ev:
		default:
			s.mu.Lock()
			s.drops++
			s.mu.Unlock()
		}
	}
}

// Topics returns the set of topics that currently have at least one subscriber.
func (b *Bus) Topics() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	out := make([]string, 0, len(b.subs))
	for t, list := range b.subs {
		if len(list) > 0 {
			out = append(out, t)
		}
	}
	return out
}

func (b *Bus) unsubscribe(s *Subscription) {
	b.mu.Lock()
	defer b.mu.Unlock()

	list := b.subs[s.topic]
	for i, sub := range list {
		if sub == s {
			b.subs[s.topic] = append(list[:i], list[i+1:]...)
			break
		}
	}
	close(s.ch)
}

// copy_slice copies src into dst and returns the number of elements copied.
func copy_slice(dst, src []*Subscription) int {
	n := len(src)
	if len(dst) < n {
		n = len(dst)
	}
	for i := 0; i < n; i++ {
		dst[i] = src[i]
	}
	return n
}
