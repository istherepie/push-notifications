package webserver

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/istherepie/push-notifications/eventbroker"
)

type Payload struct {
	Message string `json:"message"`
}

func (p *Payload) IsValid() bool {
	return p.Message != ""
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

	if err != nil || !payload.IsValid() {
		http.Error(w, "Invalid request!", http.StatusBadRequest)
		return
	}
	html.EscapeString("sdfsdf")
	// This handler should only publish messages of type "message"
	defaultType := "message"

	// Broadcast
	m.Broker.Publish(defaultType, payload.Message)
	w.WriteHeader(http.StatusNoContent)
}

type NotificationHandler struct {
	Broker *eventbroker.Broker
}

func (n *NotificationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	// Inject a heartbeat
	// into the subscription channel
	go InjectHeartbeat(r, sub)

	for {
		select {
		case <-r.Context().Done():
			log.Println("NOTIFICATIONS cLient disconnect")
			return
		case msg := <-sub.Next():
			fmt.Fprintf(w, "event: %v\n", msg.Type)
			fmt.Fprintf(w, "data: %v\n\n", msg.Value)
			flusher.Flush()
		}
	}
}

func Mux(broker *eventbroker.Broker) *http.ServeMux {

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
