package event_source

import "sync"

type Event[T any] struct {
	Type string
	Data T
}

type Broadcast[T any, K comparable] struct {
	msg map[K]chan Event[T]
	mut sync.Mutex
}

func New[T any, K comparable]() *Broadcast[T, K] {
	return &Broadcast[T, K]{
		msg: make(map[K]chan Event[T]),
		mut: sync.Mutex{},
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
	b.msg[id] = make(chan Event[T], 5)
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

	for id, _ := range b.msg {
		b.msg[id] <- e
	}
}
func (b *Broadcast[T, K]) PublishMsgExcludeID(eventType string, data T, excludeID K) {
	e := Event[T]{
		Type: eventType,
		Data: data,
	}

	b.mut.Lock()
	defer b.mut.Unlock()

	for id := range b.msg {
		if excludeID != id {
			b.msg[id] <- e
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

	for id := range b.msg {
		if !excluded[id] {
			b.msg[id] <- e
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
