package metrics

import (
	"testing"
	"time"
)

func TestActiveConnections(t *testing.T) {

	c := Counter{}

	c.Connected()
	c.Connected()
	c.Connected()

	expected := 3

	if c.Connections != expected {
		t.Error("Incorrect amount of connections\n")
		t.Logf("Expected: %d", expected)
		t.Logf("Got: %d", c.Connections)
	}
}

func TestNegativeAmountConnections(t *testing.T) {

	c := Counter{}

	c.Connected()
	c.Connected()
	c.Connected()
	c.Connected()

	c.Disconnected()
	c.Disconnected()
	c.Disconnected()
	c.Disconnected()
	c.Disconnected()
	c.Disconnected()
	c.Disconnected()
	c.Disconnected()
	c.Disconnected()

	expected := 0

	if c.Connections != expected {
		t.Error("Incorrect amount of connections\n")
		t.Logf("Expected: %d", expected)
		t.Logf("Got: %d", c.Connections)
	}
}

func TestAddVisitors(t *testing.T) {

	c := Counter{}

	c.Visitor("10.10.10.1")
	c.Visitor("10.10.10.1")
	c.Visitor("10.10.10.1")
	c.Visitor("10.10.10.1")

	result := len(c.Visitors)

	expected := 4

	if result != expected {
		t.Error("Incorrect amount of visitors\n")
		t.Logf("Expected: %d", expected)
		t.Logf("Got: %d", result)
	}
}

func TestUniqueVisitors(t *testing.T) {

	c := Counter{}

	c.Visitor("10.10.10.1")
	c.Visitor("10.10.10.1")
	c.Visitor("10.10.10.1")
	c.Visitor("10.10.10.1")
	c.Visitor("192.168.0.1")
	c.Visitor("192.168.0.1")
	c.Visitor("192.168.0.1")
	c.Visitor("172.16.100.1")
	c.Visitor("172.16.100.1")
	c.Visitor("172.16.100.1")
	c.Visitor("172.16.100.1")

	result := c.UniqueVisitors()

	expected := 3

	if result != expected {
		t.Error("Incorrect amount of visitors\n")
		t.Logf("Expected: %d", expected)
		t.Logf("Got: %d", result)
	}
}

func TestAddMessages(t *testing.T) {

	c := Counter{}

	c.Message("test")
	c.Message("test")
	c.Message("test")
	c.Message("test")

	result := len(c.Messages)

	expected := 4

	if result != expected {
		t.Error("Incorrect amount of messages\n")
		t.Logf("Expected: %d", expected)
		t.Logf("Got: %d", result)
	}
}

func TestMessagesLastHour(t *testing.T) {

	c := Counter{}

	now := time.Now()

	moreThan2Hours := now.Add(-140 * time.Minute)
	lessThan1Hours := now.Add(-45 * time.Minute)

	// Add messages manually
	c.Messages = []Message{
		Message{Type: "test", Received: moreThan2Hours},
		Message{Type: "test", Received: moreThan2Hours},
		Message{Type: "test", Received: moreThan2Hours},
		Message{Type: "test", Received: moreThan2Hours},
		Message{Type: "test", Received: lessThan1Hours}, // < 1 hours
		Message{Type: "test", Received: lessThan1Hours}, // < 1 hours
	}

	result := c.MessagesLastHour()

	expected := 2

	if result != expected {
		t.Error("Incorrect amount of messages\n")
		t.Logf("Expected: %d", expected)
		t.Logf("Got: %d", result)
	}
}
