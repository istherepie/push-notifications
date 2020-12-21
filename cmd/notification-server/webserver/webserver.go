package webserver

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/istherepie/push-notifications/eventbroker"
	"github.com/istherepie/push-notifications/metrics"
)

type Payload struct {
	Message string `json:"message"`
}

func (p *Payload) Escaped() string {
	return html.EscapeString(p.Message)
}

func (p *Payload) IsValid() bool {
	return p.Message != ""
}

type IndexHandler struct{}

func (i IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "index")
}

type Metrics struct {
	Status           string `json:"status"` // PLACEHOLDER
	OpenConnections  int    `json:"open_connections"`
	Visitors         int    `json:"visitors_total"`
	VisitorsUnique   int    `json:"visitors_unique"`
	Messages         int    `json:"messages_total"`
	MessagesLastHour int    `json:"messages_last_hour"`
}

type MetricsHandler struct {
	Counter *metrics.Counter
}

func (m *MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Set header for JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")

	metr := &Metrics{
		Status:           "ok",
		OpenConnections:  m.Counter.Connections,
		Visitors:         len(m.Counter.Visitors),
		VisitorsUnique:   m.Counter.UniqueVisitors(),
		Messages:         len(m.Counter.Messages),
		MessagesLastHour: m.Counter.MessagesLastHour(),
	}

	response, err := json.Marshal(metr)

	if err != nil {
		http.Error(w, "PARSE_RESPONSE_TO_JSON_ERR", http.StatusInternalServerError)
		return
	}

	w.Write(response)
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

	// Broadcast
	m.Broker.Publish(MessageTypeDefault, payload.Escaped())
	w.WriteHeader(http.StatusNoContent)
}

type NotificationHandler struct {
	Broker  *eventbroker.Broker
	Counter *metrics.Counter
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

	// Store visitor
	visitor := GetVisitorAddress(r)
	n.Counter.Visitor(visitor)

	// Before subscribing, annouce presence
	n.Broker.Publish(MessageTypeService, "Someone has connected!")

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

func Mux(broker *eventbroker.Broker, counter *metrics.Counter) *http.ServeMux {

	// Handlers
	handleIndex := &IndexHandler{}
	handleMessage := &MessageHandler{broker}
	handleNotifications := &NotificationHandler{broker, counter}
	handleMetrics := &MetricsHandler{counter}

	// Multiplexer
	mux := http.NewServeMux()

	// Register routes
	mux.Handle("/", handleIndex)
	mux.Handle("/metrics", handleMetrics)
	mux.Handle("/message", handleMessage)
	mux.Handle("/notifications", handleNotifications)

	return mux
}
