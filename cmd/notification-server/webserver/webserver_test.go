package webserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/istherepie/push-notifications/eventbroker"
	"github.com/istherepie/push-notifications/metrics"
)

func Setup(t *testing.T) *eventbroker.Broker {

	hook := func(status int) {
		// do nothing
	}

	// Setup Broker
	broker := &eventbroker.Broker{
		Quit:          make(chan struct{}),
		Subscriptions: make(map[*eventbroker.Subscription]struct{}),
		Register:      make(chan *eventbroker.Subscription),
		Unregister:    make(chan *eventbroker.Subscription),
		MessageQueue:  make(chan *eventbroker.Message),
		EventHook:     hook,
	}

	go broker.Run()

	t.Cleanup(func() { broker.Close() })

	return broker
}

func TestIndexHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	// Run handler
	broker := Setup(t)
	counter := &metrics.Counter{}
	handler := Mux(broker, counter)
	handler.ServeHTTP(rec, req)

	// Test
	if rec.Code != http.StatusOK {
		t.Errorf("Incorrect status code, got %v want %v", rec.Code, http.StatusOK)
	}

	expected := "index\n"
	result := rec.Body.String()

	if result != expected {
		t.Errorf("Incorrect response body, got %v want %v", result, expected)
	}
}

func TestMessageHandler(t *testing.T) {

	payload := Payload{
		Message: "test message",
	}

	data, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "/message", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	// Run handler
	broker := Setup(t)
	counter := &metrics.Counter{}
	handler := Mux(broker, counter)
	handler.ServeHTTP(rec, req)

	// Test
	if rec.Code != http.StatusNoContent {
		t.Errorf("Incorrect status code, got %v want %v", rec.Code, http.StatusOK)
	}
}

func TestMessageHandlerBadRequest(t *testing.T) {

	payload := Payload{
		Message: "",
	}

	data, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "/message", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	// Run handler
	broker := Setup(t)
	counter := &metrics.Counter{}
	handler := Mux(broker, counter)
	handler.ServeHTTP(rec, req)

	// Test
	if rec.Code != http.StatusBadRequest {
		t.Errorf("This should result in 400, got %v want %v", rec.Code, http.StatusOK)
	}

}

func TestNotificationsHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/notifications", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	// Run handler
	broker := Setup(t)
	counter := &metrics.Counter{}
	handler := Mux(broker, counter)
	go handler.ServeHTTP(rec, req)

	// Test
	if rec.Code != http.StatusOK {
		t.Errorf("Incorrect status code, got %v want %v", rec.Code, http.StatusOK)
	}
}

type StreamRecorder struct {
	*httptest.ResponseRecorder
}

func (s *StreamRecorder) Reset() {
	s.Body = new(bytes.Buffer)
	s.Flushed = false
}

func (s *StreamRecorder) WaitForFlush() {
	for !s.Flushed {
		// Do nothing
	}
}

// TODO: Rework for -race test
func TestNotificationMessages(t *testing.T) {

	req, _ := http.NewRequest("GET", "/notifications", nil)

	// Need a custom recorder
	// the response writer buffer is not flushed before the method exits
	rec := &StreamRecorder{httptest.NewRecorder()}

	broker := Setup(t)

	waitForRegistration := make(chan int)
	broker.EventHook = func(status int) {

		if status == 1 {
			waitForRegistration <- status
			return
		}
	}

	counter := &metrics.Counter{}
	mux := NotificationHandler{broker, counter}

	go mux.ServeHTTP(rec, req)

	<-waitForRegistration

	// Pass message into the loop
	broker.Publish("message", "this is a test message")

	// Wait buffer to be flushed
	rec.WaitForFlush()

	expected := "event: message\ndata: this is a test message\n\n"

	// Test
	if rec.Body.String() != expected {
		t.Error("NOTIFICATION ERROR")
		t.Log("GOT", rec.Body.String())
		t.Log("WANT", expected)
	}
}
