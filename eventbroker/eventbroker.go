package eventbroker

import "time"

type Subscription struct {
	quit     chan struct{}
	incoming chan string
}

func (c *Subscription) Close() {
	close(c.quit)
}

func (c *Subscription) Next() <-chan string {
	return c.incoming
}

type Broker struct {
	Subscriptions map[*Subscription]struct{}
	Register      chan *Subscription
	Unregister    chan *Subscription
	MessageQueue  chan string
	EventHook     func(status int)
}

func (b *Broker) WaitForClose(sub *Subscription) {
	select {
	case <-sub.quit:
		b.Unregister <- sub
		return
	}
}

func (b *Broker) Publish(message string) {
	b.MessageQueue <- message
}

func (b *Broker) Subscribe() *Subscription {

	sub := &Subscription{
		quit:     make(chan struct{}),
		incoming: make(chan string),
	}

	b.Register <- sub
	return sub
}

func (b *Broker) Broadcast(message string) {

	transmit := func(sub *Subscription, message string) {
		select {
		case sub.incoming <- message:
			return
		// TODO: The timeout should be configurable
		case <-time.After(1500 * time.Millisecond):
			return
		}
	}

	for sub := range b.Subscriptions {
		go transmit(sub, message)
	}
}

func (b *Broker) Run() {
	for {
		select {
		case sub := <-b.Register:
			var empty struct{}
			b.Subscriptions[sub] = empty
			b.EventHook(1)
		case sub := <-b.Unregister:
			delete(b.Subscriptions, sub)
			b.EventHook(2)
		case message := <-b.MessageQueue:
			b.Broadcast(message)
			b.EventHook(3)
		}
	}
}
