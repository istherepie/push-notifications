package main

import (
	"log"
	"net"
	"net/http"

	"github.com/istherepie/push-notifications/cmd/notification-server/webserver"
	"github.com/istherepie/push-notifications/eventbroker"
	"github.com/istherepie/push-notifications/metrics"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")

	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	// Metrics
	counter := &metrics.Counter{}

	recordMessages := func(status int) {
		switch status {
		case 1:
			counter.Connected()
			return
		case 2:
			counter.Disconnected()
			return
		case 3:
			// Cannot get message type from the eventhook
			// For now, just set the type to `any`
			// in order to record the total sum
			counter.Message("any")
			return
		}
	}

	broker := &eventbroker.Broker{
		Subscriptions: make(map[*eventbroker.Subscription]struct{}),
		Register:      make(chan *eventbroker.Subscription),
		Unregister:    make(chan *eventbroker.Subscription),
		MessageQueue:  make(chan *eventbroker.Message),
		EventHook:     recordMessages,
	}

	go broker.Run()

	mux := webserver.Mux(broker, counter)

	serveErr := http.Serve(listener, mux)

	if serveErr != nil {
		log.Fatal(serveErr)
	}
}
