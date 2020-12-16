package webserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/istherepie/push-notifications/eventbroker"
)

type Payload struct {
	Message string `json:"message"`
}

type IndexHandler struct{}

func (i IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "index")
}

type MessageHandler struct {
	Broker *eventbroker.Broker
}

func (m MessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload Payload

	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		http.Error(w, "Invalid request!", http.StatusBadRequest)
		return
	}

	// Broadcast
	m.Broker.Publish(payload.Message)
	w.WriteHeader(http.StatusNoContent)
}

type NotificationHandler struct {
	Broker *eventbroker.Broker
}

func (n NotificationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("NOTIFICATIONS cLient connected")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)

	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	sub := n.Broker.Subscribe()
	defer sub.Close()

	for {
		select {
		case <-r.Context().Done():
			log.Println("NOTIFICATIONS cLient disconnect")
			return
		case msg := <-sub.Next():
			fmt.Fprintf(w, "data: %v\n\n", msg)
			flusher.Flush()
		}
	}
}

func Mux() *http.ServeMux {

	// EventBroker
	broker := &eventbroker.Broker{}
	broker.Init()

	// Handlers
	handleIndex := &IndexHandler{}
	handleMessage := &MessageHandler{broker}
	handleNotifications := &NotificationHandler{broker}

	// Multiplexer
	mux := http.NewServeMux()

	// Register routes
	mux.Handle("/", handleIndex)
	mux.Handle("/message", handleMessage)
	mux.Handle("/notifications", handleNotifications)

	return mux
}
