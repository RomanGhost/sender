package message_test

import (
	"sender/internal/server/blockchain/p2pprotocol/message"
	"testing"
	"time"
)

func TestBaseMessage_GetSetID(t *testing.T) {
	msg := &message.BaseMessage{}

	// Проверяем SetID
	expectedID := uint64(42)
	msg.SetID(expectedID)

	if msg.GetID() != expectedID {
		t.Errorf("Expected ID to be %d, got %d", expectedID, msg.GetID())
	}
}

func TestBaseMessage_Timestamp(t *testing.T) {
	now := time.Now()
	msg := &message.BaseMessage{
		Timestamp: now,
	}

	// Проверяем Timestamp
	if !msg.Timestamp.Equal(now) {
		t.Errorf("Expected Timestamp to be %v, got %v", now, msg.Timestamp)
	}
}

func TestBaseMessage_MessageType(t *testing.T) {
	msg := &message.BaseMessage{}

	// Проверяем MessageType
	expectedType := "BaseMessage"
	if msg.MessageType() != expectedType {
		t.Errorf("Expected MessageType to be '%s', got '%s'", expectedType, msg.MessageType())
	}
}

func TestMessageType_String(t *testing.T) {
	tests := []struct {
		messageType message.MessageType
		expected    string
	}{
		{message.RequestBlocksBeforeMessage, "RequestBlocksBeforeMessage"},
		{message.RequestMessageInfo, "RequestMessageInfo"},
		{message.RequestLastNBlocksMessage, "RequestLastNBlocksMessage"},
		{message.ResponseBlockMessage, "ResponseBlockMessage"},
		{message.ResponseChainMessage, "ResponseChainMessage"},
		{message.ResponsePeerMessage, "ResponsePeerMessage"},
		{message.ResponseMessageInfo, "ResponseMessageInfo"},
		{message.ResponseTransactionMessage, "ResponseTransactionMessage"},
	}

	for _, tt := range tests {
		t.Run(string(tt.messageType), func(t *testing.T) {
			if tt.messageType.String() != tt.expected {
				t.Errorf("Expected String() to return '%s', got '%s'", tt.expected, tt.messageType.String())
			}
		})
	}
}
