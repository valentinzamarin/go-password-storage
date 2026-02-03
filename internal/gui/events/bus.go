package events

import (
	"sync"
)

type Event struct {
	Name    string
	Payload interface{}
}

type Bus struct {
	mu   sync.RWMutex
	subs map[string]map[chan Event]struct{}
}

func New() *Bus {
	return &Bus{
		subs: make(map[string]map[chan Event]struct{}),
	}
}

var Default = New()

func (b *Bus) Subscribe(name string) (<-chan Event, func()) {
	ch := make(chan Event, 1)
	b.mu.Lock()
	m, ok := b.subs[name]
	if !ok {
		m = make(map[chan Event]struct{})
		b.subs[name] = m
	}
	m[ch] = struct{}{}
	b.mu.Unlock()

	unsub := func() {
		b.mu.Lock()
		if m, ok := b.subs[name]; ok {
			delete(m, ch)
			close(ch)
			if len(m) == 0 {
				delete(b.subs, name)
			}
		}
		b.mu.Unlock()
	}

	return ch, unsub
}

func (b *Bus) Publish(name string, payload interface{}) {
	b.mu.RLock()
	m := b.subs[name]
	for ch := range m {
		select {
		case ch <- Event{Name: name, Payload: payload}:
		default:
		}
	}
	b.mu.RUnlock()
}

func Subscribe(name string) (<-chan Event, func()) { return Default.Subscribe(name) }
func Publish(name string, payload interface{})     { Default.Publish(name, payload) }
