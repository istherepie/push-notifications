package main

import (
	"log"
	"net"
	"net/http"

	"github.com/istherepie/push-notifications/cmd/notification-server/webserver"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")

	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	mux := webserver.Mux()

	serveErr := http.Serve(listener, mux)

	if serveErr != nil {
		log.Fatal(serveErr)
	}
}
