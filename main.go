package event_source

import "sync"

type Event[T any] struct {
	Type string
	Data T
}

type Broadcast[T any, K comparable] struct {
	msg      map[K]chan Event[T]
	mut      sync.Mutex
	capacity int
}

func New[T any, K comparable](capacity ...int) *Broadcast[T, K] {
	c := 10
	if len(capacity) > 0 {
		c = capacity[0]
	}
	return &Broadcast[T, K]{
		msg:      make(map[K]chan Event[T]),
		mut:      sync.Mutex{},
		capacity: c,
	}
}

func (b *Broadcast[T, K]) PublishMsgToID(eventType string, data T, id K) {
	e := Event[T]{
		Type: eventType,
		Data: data,
	}

	b.mut.Lock()
	defer b.mut.Unlock()

	if ch, ok := b.msg[id]; ok {
		select {
		case ch <- e:
		default:
		}
	}
}

func (b *Broadcast[T, K]) Subscribe(id K) {
	var zero K
	if id == zero {
		return
	}
	_, ok := b.msg[id]
	if ok {
		return
	}
	b.mut.Lock()
	defer b.mut.Unlock()
	b.msg[id] = make(chan Event[T], b.capacity)
}

func (b *Broadcast[T, K]) CancelSubscribe(id K) {
	var zero K
	if id == zero {
		return
	}
	_, ok := b.msg[id]
	if !ok {
		return
	}
	b.mut.Lock()
	defer b.mut.Unlock()
	delete(b.msg, id)
}

func (b *Broadcast[T, K]) PublishMsg(eventType string, data T) {
	e := Event[T]{
		Type: eventType,
		Data: data,
	}

	b.mut.Lock()
	defer b.mut.Unlock()

	for id, ch := range b.msg {
		select {
		case ch <- e:
		default:
		}
	}
}
func (b *Broadcast[T, K]) PublishMsgExcludeID(eventType string, data T, excludeID K) {
	e := Event[T]{
		Type: eventType,
		Data: data,
	}

	b.mut.Lock()
	defer b.mut.Unlock()

	for id, ch := range b.msg {
		if excludeID != id {
			select {
			case ch <- e:
			default:
			}
		}
	}
}
func (b *Broadcast[T, K]) PublishMsgExcludeIDs(eventType string, data T, excludedIDs []K) {
	e := Event[T]{
		Type: eventType,
		Data: data,
	}

	b.mut.Lock()
	defer b.mut.Unlock()

	excluded := make(map[K]bool)
	for _, id := range excludedIDs {
		excluded[id] = true
	}

	for id, ch := range b.msg {
		if !excluded[id] {
			select {
			case ch <- e:
			default:
			}
		}
	}
}

func (b *Broadcast[T, K]) ReceiveMsg(id K) chan Event[T] {
	var zero K
	if id == zero {
		return nil
	}
	_, ok := b.msg[id]
	if !ok {
		return nil
	}
	return b.msg[id]
}
