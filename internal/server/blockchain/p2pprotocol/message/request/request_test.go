package request_test

import (
	"sender/internal/server/blockchain/p2pprotocol/message"
	"sender/internal/server/blockchain/p2pprotocol/message/request"
	"testing"
	"time"
)

func TestBaseMessage(t *testing.T) {
	baseMsg := &message.BaseMessage{
		ID:        123,
		Timestamp: time.Now(),
	}

	// Test GetID
	if baseMsg.GetID() != 123 {
		t.Errorf("Expected ID to be 123, got %d", baseMsg.GetID())
	}

	// Test SetID
	baseMsg.SetID(456)
	if baseMsg.GetID() != 456 {
		t.Errorf("Expected ID to be 456 after SetID, got %d", baseMsg.GetID())
	}

	// Test MessageType
	if baseMsg.MessageType() != "BaseMessage" {
		t.Errorf("Expected MessageType to be 'BaseMessage', got '%s'", baseMsg.MessageType())
	}
}

func TestBlocksBeforeMessage(t *testing.T) {
	msg := request.NewBlocksBeforeMessage()

	// Test ID initialization
	if msg.GetID() != 0 {
		t.Errorf("Expected initial ID to be 0, got %d", msg.GetID())
	}

	// Test Timestamp initialization
	if time.Since(msg.Timestamp) > time.Second {
		t.Errorf("Timestamp is not recent: %v", msg.Timestamp)
	}

	// Test MessageType
	expectedType := message.RequestBlocksBeforeMessage.String()
	if msg.MessageType() != expectedType {
		t.Errorf("Expected MessageType to be '%s', got '%s'", expectedType, msg.MessageType())
	}
}

func TestInfoMessage(t *testing.T) {
	msg := request.NewInfoMessage()

	// Test ID initialization
	if msg.GetID() != 0 {
		t.Errorf("Expected initial ID to be 0, got %d", msg.GetID())
	}

	// Test Timestamp initialization
	if time.Since(msg.Timestamp) > time.Second {
		t.Errorf("Timestamp is not recent: %v", msg.Timestamp)
	}

	// Test MessageType
	expectedType := message.RequestMessageInfo.String()
	if msg.MessageType() != expectedType {
		t.Errorf("Expected MessageType to be '%s', got '%s'", expectedType, msg.MessageType())
	}
}

func TestLastNBlocksMessage(t *testing.T) {
	countBlocks := uint(10)
	msg := request.NewLastNBlocksMessage(countBlocks)

	// Test ID initialization
	if msg.GetID() != 0 {
		t.Errorf("Expected initial ID to be 0, got %d", msg.GetID())
	}

	// Test Timestamp initialization
	if time.Since(msg.Timestamp) > time.Second {
		t.Errorf("Timestamp is not recent: %v", msg.Timestamp)
	}

	// Test N field
	if msg.N != countBlocks {
		t.Errorf("Expected N to be %d, got %d", countBlocks, msg.N)
	}

	// Test MessageType
	expectedType := message.RequestLastNBlocksMessage.String()
	if msg.MessageType() != expectedType {
		t.Errorf("Expected MessageType to be '%s', got '%s'", expectedType, msg.MessageType())
	}
}
