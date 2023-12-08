package event_source

import "sync"

type Event[T any] struct {
	Type string
	Data T
}

type Broadcast[T any] struct {
	msg map[string]chan Event[T]
	mut sync.Mutex
}

func New[T any]() *Broadcast[T] {
	return &Broadcast[T]{
		msg: make(map[string]chan Event[T]),
		mut: sync.Mutex{},
	}
}

func (b *Broadcast[T]) Subscribe(id string) {
	if id == "" {
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

func (b *Broadcast[T]) CancelSubscribe(id string) {
	if id == "" {
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

func (b *Broadcast[T]) PublishMsg(e Event[T]) {
	b.mut.Lock()
	defer b.mut.Unlock()
	for id, _ := range b.msg {
		b.msg[id] <- e
	}
}

func (b *Broadcast[T]) ReceiveMsg(id string) chan Event[T] {
	if id == "" {
		return nil
	}
	_, ok := b.msg[id]
	if !ok {
		return nil
	}
	return b.msg[id]
}
