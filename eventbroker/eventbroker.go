package event

import (
	"sync"
	"time"
)

type Subscription struct {
	listener chan string
	quit     chan struct{}
}

func (s *Subscription) Next() <-chan string {
	return s.listener
}

func (s *Subscription) Close() {
	close(s.quit)
}

type Broker struct {
	subscriptions map[*Subscription]struct{}
	mtx           sync.RWMutex
}

func (b *Broker) Init() {
	b.subscriptions = make(map[*Subscription]struct{})
}

func (b *Broker) addSub(sub *Subscription) {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	var empty struct{}
	b.subscriptions[sub] = empty
}

func (b *Broker) remSub(sub *Subscription) {
	b.mtx.Lock()
	defer b.mtx.Unlock()
	close(sub.listener)
	delete(b.subscriptions, sub)
}

func (b *Broker) waitForClose(sub *Subscription) {
	<-sub.quit
	b.remSub(sub)
}

func (b *Broker) Subscribe() *Subscription {

	subscription := &Subscription{
		listener: make(chan string, 1),
		quit:     make(chan struct{}),
	}

	b.addSub(subscription)

	// Remove when closed
	go b.waitForClose(subscription)

	return subscription
}

func (b *Broker) GetSubs() []*Subscription {
	b.mtx.RLock()
	defer b.mtx.RUnlock()

	var subscribers []*Subscription

	for sub := range b.subscriptions {
		subscribers = append(subscribers, sub)
	}

	return subscribers
}

func (b *Broker) CountSubs() int {
	b.mtx.RLock()
	defer b.mtx.RUnlock()
	return len(b.subscriptions)
}

func (b *Broker) Publish(message string) {

	if b.CountSubs() == 0 {
		return
	}

	broadcast := func(sub *Subscription, message string) {
		select {
		case sub.listener <- message:
			return

		// TODO: The timeout should be configurable
		case <-time.After(300 * time.Millisecond):
			return
		}
	}

	for _, subscription := range b.GetSubs() {
		go broadcast(subscription, message)
	}

}
