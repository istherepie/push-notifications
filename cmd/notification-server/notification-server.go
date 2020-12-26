package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/istherepie/push-notifications/cmd/notification-server/webserver"
	"github.com/istherepie/push-notifications/eventbroker"
	"github.com/istherepie/push-notifications/metrics"
)

func main() {
	hostname := flag.String("host", "localhost", "Set the host address of the service.")
	port := flag.Int("port", 8080, "Set the tcp port of the service.")
	static := flag.String("static", "", "Specify a directory to serve static files from.")
	flag.Parse()

	var webDir string

	if *static != "" {
		path := fmt.Sprintf("%v/index.html", *static)
		_, err := os.Stat(path)

		if os.IsNotExist(err) {
			log.Printf("No index.html found in path: %v", *static)
		}

		webDir = *static
	}

	addr := fmt.Sprintf("%v:%d", *hostname, *port)
	listener, err := net.Listen("tcp", addr)

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

	mux := webserver.Mux(webDir, broker, counter)

	serveErr := http.Serve(listener, mux)

	if serveErr != nil {
		log.Fatal(serveErr)
	}
}
