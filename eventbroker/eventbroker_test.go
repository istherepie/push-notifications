package eventbroker

import (
	"testing"
)

func TestSubscriberCount(t *testing.T) {
	broker := Broker{}
	broker.Init()

	broker.Subscribe()

	// Count subs
	result := broker.CountSubs()

	// Expecting 1 sub
	expected := 1

	if result != expected {
		t.Errorf("Incorrect amount of subs - result: %d (Should be %d)", result, expected)
	}
}

func TestAddAndRemoveSubs(t *testing.T) {
	broker := Broker{}
	broker.Init()

	// Create 4 subs
	sub1 := broker.Subscribe()
	sub2 := broker.Subscribe()
	sub3 := broker.Subscribe()
	sub4 := broker.Subscribe()

	// Close subs
	sub1.Close()
	sub2.Close()
	sub3.Close()
	sub4.Close()

	// Wait for all to close
	<-sub1.Next()
	<-sub2.Next()
	<-sub3.Next()
	<-sub4.Next()

	// Count subs
	result := broker.CountSubs()

	// Expecting 1 sub
	expected := 0

	if result != expected {
		t.Errorf("Incorrect amount of subs - result: %d (Should be %d)", result, expected)
	}
}

func TestEventNotification(t *testing.T) {

	broker := Broker{}
	broker.Init()

	sub := broker.Subscribe()

	// Fixtures
	testmessages := []string{
		"test message 1",
		"test message 2",
		"test message 3",
		"test message 4",
		"test message 5",
	}

	// Test all 5 messages
	for i := 0; i < 5; i++ {
		broker.Publish(testmessages[i])

		// Get message
		result := <-sub.Next()

		if result != testmessages[i] {
			t.Errorf("Message I/O error, expected: %v , got: %v", result, testmessages[i])
		}

	}

}
