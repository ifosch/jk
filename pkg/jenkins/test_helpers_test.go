package jenkins

import "testing"

func assertExpectedMessages(expectedMsgs []Message, ch chan Message, t *testing.T) {
	for _, expectedMsg := range expectedMsgs {
		msg := <-ch
		if msg.Message != expectedMsg.Message {
			t.Fatalf("Unexpected Message in reply: %s != %s", msg.Message, expectedMsg.Message)
		}
		if msg.Error != expectedMsg.Error {
			t.Fatalf("Unexpected Error in reply: %v != %v", msg.Error, expectedMsg.Error)
		}
		if msg.Done != expectedMsg.Done {
			t.Fatalf("Unexpected Done in reply: %v != %v", msg.Done, expectedMsg.Done)
		}
	}
}
