package eventbroker

import (
	"sync"
	"time"
)

type Message struct {
	Type  string
	Value string
}

type Subscription struct {
	quit     chan struct{}
	incoming chan *Message
}

func (c *Subscription) Close() {
	close(c.quit)
}

func (c *Subscription) Next() <-chan *Message {
	return c.incoming
}

type Broker struct {
	mtx           sync.RWMutex
	Quit          chan struct{}
	Subscriptions map[*Subscription]struct{}
	Register      chan *Subscription
	Unregister    chan *Subscription
	MessageQueue  chan *Message
	EventHook     func(status int)
}

func (b *Broker) WaitForClose(sub *Subscription) {
	select {
	case <-sub.quit:
		b.Unregister <- sub
		return
	}
}

func (b *Broker) Publish(msgType string, msgValue string) {
	b.MessageQueue <- &Message{msgType, msgValue}
}

func (b *Broker) Subscribe() *Subscription {

	sub := &Subscription{
		quit:     make(chan struct{}),
		incoming: make(chan *Message),
	}

	b.Register <- sub

	return sub
}

func (b *Broker) Broadcast(msg *Message) {

	transmit := func(sub *Subscription, msg *Message) {
		select {
		case sub.incoming <- msg:
			return
		// TODO: The timeout should be configurable
		case <-time.After(1500 * time.Millisecond):
			return
		}
	}

	b.mtx.RLock()
	for sub := range b.Subscriptions {
		go transmit(sub, msg)
	}

	b.mtx.RUnlock()
}

func (b *Broker) CountSubs() int {
	b.mtx.RLock()
	defer b.mtx.RUnlock()
	return len(b.Subscriptions)
}

func (b *Broker) Close() {
	var empty struct{}
	b.Quit <- empty
}

func (b *Broker) Run() {
	for {
		select {
		case sub := <-b.Register:
			b.mtx.Lock()
			var empty struct{}
			b.Subscriptions[sub] = empty
			b.mtx.Unlock()
			go b.WaitForClose(sub)
			b.EventHook(1)
		case sub := <-b.Unregister:
			close(sub.incoming)
			b.mtx.Lock()
			delete(b.Subscriptions, sub)
			b.mtx.Unlock()
			b.EventHook(2)
		case message := <-b.MessageQueue:
			b.Broadcast(message)
			b.EventHook(3)
		case <-b.Quit:
			return
		}
	}
}
