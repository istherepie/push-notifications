package main

import (
	"log"
	"net"
	"net/http"

	"github.com/istherepie/push-notifications/cmd/notification-server/webserver"
	"github.com/istherepie/push-notifications/eventbroker"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")

	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	broker := &eventbroker.Broker{
		Subscriptions: make(map[*eventbroker.Subscription]struct{}),
		Register:      make(chan *eventbroker.Subscription),
		Unregister:    make(chan *eventbroker.Subscription),
		MessageQueue:  make(chan *eventbroker.Message),
		EventHook:     func(status int) {},
	}

	go broker.Run()

	mux := webserver.Mux(broker)

	serveErr := http.Serve(listener, mux)

	if serveErr != nil {
		log.Fatal(serveErr)
	}
}
