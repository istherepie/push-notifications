package metrics

import (
	"sync"
	"time"
)

type Message struct {
	Type     string
	Received time.Time
}

type Visitor struct {
	Address string
}

type Counter struct {
	mtx         sync.RWMutex
	Messages    []Message
	Visitors    []Visitor
	Connections int
}

func (c *Counter) Connected() {
	c.Connections++
}

func (c *Counter) Disconnected() {
	if c.Connections <= 0 {
		return
	}

	c.Connections--
}

func (c *Counter) Visitor(address string) {
	visitor := Visitor{address}

	c.Visitors = append(c.Visitors, visitor)
}

func (c *Counter) Message(msgType string) {
	message := Message{
		Type:     msgType,
		Received: time.Now(),
	}

	c.Messages = append(c.Messages, message)
}

func (c *Counter) UniqueVisitors() int {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	m := make(map[string]struct{})

	for _, visitor := range c.Visitors {
		var empty struct{}
		m[visitor.Address] = empty
	}
	return len(m)
}

func (c *Counter) MessagesLastHour() int {
	var counter int

	for _, message := range c.Messages {
		now := time.Now()
		timeDiff := now.Sub(message.Received)

		if timeDiff > time.Hour {
			continue
		}

		counter++
	}

	return counter
}
