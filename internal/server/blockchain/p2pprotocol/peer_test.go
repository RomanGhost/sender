package p2pprotocol_test

import (
	"net"
	"sender/internal/server/blockchain/p2pprotocol"
	"sender/internal/server/blockchain/p2pprotocol/message"
	"sender/internal/server/blockchain/p2pprotocol/message/request"
	"sync"
	"testing"
)

type MockConnectionPool struct {
	mu         sync.Mutex
	peers      map[string]string
	broadcasts []string
}

func NewMockConnectionPool() *MockConnectionPool {
	return &MockConnectionPool{
		peers:      make(map[string]string),
		broadcasts: []string{},
	}
}

func (m *MockConnectionPool) AddPeer(address string, conn net.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.peers[address] = address
}

func (m *MockConnectionPool) RemovePeer(address string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.peers, address)
}

func (m *MockConnectionPool) GetPeerAddresses() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	addresses := []string{}
	for addr := range m.peers {
		addresses = append(addresses, addr)
	}
	return addresses
}

func (m *MockConnectionPool) Broadcast(message string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.broadcasts = append(m.broadcasts, message)
}

func (m *MockConnectionPool) GetBroadcasts() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.broadcasts
}

func (m *MockConnectionPool) ClearBroadcasts() {
	m.mu.Lock()
	defer m.mu.Unlock()
	clear(m.broadcasts)
}

func TestP2PProtocol_HandleMessage(t *testing.T) {

}

func TestP2PProtocol_Broadcast(t *testing.T) {
	mockPool := NewMockConnectionPool()
	sender := make(chan message.Message, 1)
	protocol := p2pprotocol.New(mockPool, sender)

	msg := request.NewInfoMessage()
	protocol.Broadcast(msg, false)

	if len(mockPool.GetBroadcasts()) != 1 {
		t.Fatalf("Expected 1 broadcast, got %d", len(mockPool.GetBroadcasts()))
	}

	if msg.GetID() != protocol.GetLastMessageId() {
		t.Errorf("Expected message ID to match protocol's lastMessageID")
	}
}

func TestP2PProtocol_HandleDisconnectedPeers(t *testing.T) {
	mockPool := NewMockConnectionPool()
	sender := make(chan message.Message, 1)
	protocol := p2pprotocol.New(mockPool, sender)

	// Добавляем пиров
	mockPool.AddPeer("peer1", nil)
	mockPool.AddPeer("peer2", nil)

	// Удаляем одного из пиров
	mockPool.RemovePeer("peer2")

	msg := request.NewInfoMessage()
	protocol.Broadcast(msg, false)

	if len(mockPool.GetBroadcasts()) != 1 {
		t.Fatalf("Expected 1 broadcast, got %d", len(mockPool.GetBroadcasts()))
	}

	addresses := mockPool.GetPeerAddresses()
	if len(addresses) != 1 || addresses[0] != "peer1" {
		t.Errorf("Expected only peer1 to be connected, got: %v", addresses)
	}
}
