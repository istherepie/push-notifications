package webserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndexHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	// Run handler
	handler := Mux()
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
	handler := Mux()
	handler.ServeHTTP(rec, req)

	// Test
	if rec.Code != http.StatusNoContent {
		t.Errorf("Incorrect status code, got %v want %v", rec.Code, http.StatusOK)
	}
}

func TestNotificationsHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/notifications", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()

	// Run handler
	handler := Mux()
	go handler.ServeHTTP(rec, req)

	// Test
	if rec.Code != http.StatusOK {
		t.Errorf("Incorrect status code, got %v want %v", rec.Code, http.StatusOK)
	}
}
