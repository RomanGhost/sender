package responce_test

import (
	"sender/internal/data/blockchain/block"
	"sender/internal/data/blockchain/transaction"
	"sender/internal/server/blockchain/p2pprotocol/message"
	"sender/internal/server/blockchain/p2pprotocol/message/responce"
	"testing"
	"time"
)

func TestChainMessage(t *testing.T) {
	blocks := []block.Block{
		{ID: 1, TimeCreated: time.Now()},
		{ID: 2, TimeCreated: time.Now()},
	}
	msg := responce.NewChainMessage(blocks)

	// Test MessageType
	expectedType := message.ResponseChainMessage.String()
	if msg.MessageType() != expectedType {
		t.Errorf("Expected MessageType to be '%s', got '%s'", expectedType, msg.MessageType())
	}

	// Test ChainMessage field
	if len(msg.ChainMessage) != len(blocks) {
		t.Errorf("Expected ChainMessage to have %d blocks, got %d", len(blocks), len(msg.ChainMessage))
	}
}

func TestBlockMessage(t *testing.T) {
	testBlock := &block.Block{ID: 123, TimeCreated: time.Now()}
	force := true
	msg := responce.NewBlockMessage(testBlock, force)

	// Test MessageType
	expectedType := message.ResponseBlockMessage.String()
	if msg.MessageType() != expectedType {
		t.Errorf("Expected MessageType to be '%s', got '%s'", expectedType, msg.MessageType())
	}

	// Test Block and Force fields
	if msg.Block != testBlock {
		t.Errorf("Expected Block to be %+v, got %+v", testBlock, msg.Block)
	}
	if msg.Force != force {
		t.Errorf("Expected Force to be %v, got %v", force, msg.Force)
	}
}

func TestInfoMessage(t *testing.T) {
	msg := responce.NewInfoMessage()

	// Test MessageType
	expectedType := message.ResponseMessageInfo.String()
	if msg.MessageType() != expectedType {
		t.Errorf("Expected MessageType to be '%s', got '%s'", expectedType, msg.MessageType())
	}
}

func TestPeerMessage(t *testing.T) {
	peers := []string{"peer1:8080", "peer2:9090"}
	msg := responce.NewPeerMessage(peers)

	// Test MessageType
	expectedType := message.ResponsePeerMessage.String()
	if msg.MessageType() != expectedType {
		t.Errorf("Expected MessageType to be '%s', got '%s'", expectedType, msg.MessageType())
	}

	// Test PeerAddresses field
	if len(msg.PeerAddresses) != len(peers) {
		t.Errorf("Expected PeerAddresses to have %d items, got %d", len(peers), len(msg.PeerAddresses))
	}
	for i, peer := range peers {
		if msg.PeerAddresses[i] != peer {
			t.Errorf("Expected PeerAddresses[%d] to be '%s', got '%s'", i, peer, msg.PeerAddresses[i])
		}
	}
}

func TestTransactionMessage(t *testing.T) {
	testTransaction := &transaction.Transaction{
		Sender:   "sender1",
		Transfer: 100.50,
	}
	msg := responce.NewTransactionMessage(testTransaction)

	// Test MessageType
	expectedType := message.ResponseTransactionMessage.String()
	if msg.MessageType() != expectedType {
		t.Errorf("Expected MessageType to be '%s', got '%s'", expectedType, msg.MessageType())
	}

	// Test Transaction field
	if msg.Transaction != testTransaction {
		t.Errorf("Expected Transaction to be %+v, got %+v", testTransaction, msg.Transaction)
	}
}
