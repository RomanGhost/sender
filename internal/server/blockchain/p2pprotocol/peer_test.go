package p2pprotocol_test

import (
	"fmt"
	"net"
	"sender/internal/data/blockchain/block"
	"sender/internal/data/blockchain/transaction"
	"sender/internal/server/blockchain/p2pprotocol"
	"sender/internal/server/blockchain/p2pprotocol/message"
	"sender/internal/server/blockchain/p2pprotocol/message/request"
	"sender/internal/server/blockchain/p2pprotocol/message/responce"
	"sender/internal/server/blockchain/p2pprotocol/serializemessage"
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

func TestP2PProtocol_HandleMessage_RequestInfoMessage(t *testing.T) {
	connextion_pool := NewMockConnectionPool()
	channel := make(chan message.Message)

	protocol := p2pprotocol.New(connextion_pool, channel)

	infoMessage := request.NewInfoMessage()
	infoMessageSerialize := serializemessage.NewGenericMessage(infoMessage)
	jsonInfoMessageSerialize, _ := infoMessageSerialize.ToJSON()
	go func() {
		protocol.HandleMessage(string(jsonInfoMessageSerialize))
		defer close(channel)
	}()

	channelRes := <-channel
	for i := range channel {
		fmt.Println(i)
	}
	fmt.Println(connextion_pool.GetBroadcasts()[0], channelRes)
	if channelRes != nil {
		t.Fatalf("Expected no message in channel, got %+v", channelRes)
	}
	messageFromJson, _ := serializemessage.FromJSON([]byte(connextion_pool.GetBroadcasts()[0]))
	if messageFromJson.Content.MessageType() != message.ResponseMessageInfo.String() {
		t.Fatalf("Expected message type %v, got %+v", message.ResponseMessageInfo.String(), messageFromJson.Content.MessageType())
	}
}

func TestP2PProtocol_HandleMessage_ResponseInfoMessage(t *testing.T) {
	connextion_pool := NewMockConnectionPool()
	channel := make(chan message.Message)

	protocol := p2pprotocol.New(connextion_pool, channel)

	infoMessage := responce.NewInfoMessage()
	infoMessage.SetID(5)
	infoMessageSerialize := serializemessage.NewGenericMessage(infoMessage)
	jsonInfoMessageSerialize, _ := infoMessageSerialize.ToJSON()
	go func() {
		protocol.HandleMessage(string(jsonInfoMessageSerialize))
		defer close(channel)
	}()

	channelRes := <-channel
	if channelRes != nil {
		t.Fatalf("Expected no message in channel, got %+v", channelRes)
	}
	if len(connextion_pool.GetBroadcasts()) != 0 {
		t.Fatalf("Expected no message, got %+v", connextion_pool.GetBroadcasts())
	}
	if protocol.GetLastMessageId() != 5 {
		t.Fatalf("Expected message id, got %d", protocol.GetLastMessageId())
	}
}

func TestP2PProtocol_HandleMessage_ResponseInfoMessageLessID(t *testing.T) {
	connextion_pool := NewMockConnectionPool()
	channel := make(chan message.Message)

	protocol := p2pprotocol.New(connextion_pool, channel)

	infoMessage := responce.NewInfoMessage()
	infoMessage.SetID(5)
	infoMessageSerialize := serializemessage.NewGenericMessage(infoMessage)
	jsonInfoMessageSerialize, _ := infoMessageSerialize.ToJSON()
	protocol.HandleMessage(string(jsonInfoMessageSerialize))

	//Уменьшаем messageID
	infoMessage.SetID(1)
	infoMessageSerialize = serializemessage.NewGenericMessage(infoMessage)
	jsonInfoMessageSerialize, _ = infoMessageSerialize.ToJSON()
	go func() {
		protocol.HandleMessage(string(jsonInfoMessageSerialize))
		defer close(channel)
	}()

	channelRes := <-channel
	if channelRes != nil {
		t.Fatalf("Expected no message in channel, got %+v", channelRes)
	}
	if len(connextion_pool.GetBroadcasts()) != 0 {
		t.Fatalf("Expected no message, got %+v", connextion_pool.GetBroadcasts())
	}
	if protocol.GetLastMessageId() != 5 {
		t.Fatalf("Expected message id, got %d", protocol.GetLastMessageId())
	}
}

func TestP2PProtocol_HandleMessage_OtherMessage(t *testing.T) {
	connextion_pool := NewMockConnectionPool()
	channel := make(chan message.Message, 1) // Используем буферизированный канал

	protocol := p2pprotocol.New(connextion_pool, channel)

	infoMessage := responce.NewChainMessage([]block.Block{})
	infoMessage.SetID(1)
	infoMessageSerialize := serializemessage.NewGenericMessage(infoMessage)
	jsonInfoMessageSerialize, _ := infoMessageSerialize.ToJSON()

	go func() {
		protocol.HandleMessage(string(jsonInfoMessageSerialize))
		defer close(channel)
	}()

	// Ожидаем сообщение из канала
	channelRes := <-channel
	if channelRes == nil {
		t.Fatalf("Expected message in channel, got null")
	}

	if len(connextion_pool.GetBroadcasts()) == 0 {
		t.Fatalf("Expected message, but didn't got")
	}
}

func TestP2PProtocol_HandleMessage_UnknownMessage(t *testing.T) {
	connextion_pool := NewMockConnectionPool()
	channel := make(chan message.Message)

	protocol := p2pprotocol.New(connextion_pool, channel)

	invalidMessage := `{"Type":"UnknownType","Content":{}}`
	err := protocol.HandleMessage(invalidMessage)

	if err == nil {
		t.Fatalf("Expect error becouse UncnownType")
	}

	if len(connextion_pool.GetBroadcasts()) != 0 {
		t.Fatalf("Expected no broadcasts for unknown message type, got %+v", connextion_pool.GetBroadcasts())
	}
}

func TestP2PProtocol_RequestInfoMessage(t *testing.T) {
	mockPool := NewMockConnectionPool()
	sender := make(chan message.Message, 1)
	protocol := p2pprotocol.New(mockPool, sender)

	protocol.RequestInfoMessage()

	if len(mockPool.GetBroadcasts()) != 1 {
		t.Fatalf("Expected 1 broadcast for RequestInfoMessage, got %d", len(mockPool.GetBroadcasts()))
	}

	broadcast := mockPool.GetBroadcasts()[0]
	msg, err := serializemessage.FromJSON([]byte(broadcast))
	if err != nil {
		t.Fatalf("Failed to deserialize broadcast message: %v", err)
	}
	if msg.Content.MessageType() != message.RequestMessageInfo.String() {
		t.Fatalf("Expected RequestMessageInfo, got %v", msg.Content.MessageType())
	}
}

func TestP2PProtocol_RequestChainMessage(t *testing.T) {
	mockPool := NewMockConnectionPool()
	sender := make(chan message.Message, 1)
	protocol := p2pprotocol.New(mockPool, sender)

	protocol.RequestChainMessage(5)

	if len(mockPool.GetBroadcasts()) != 1 {
		t.Fatalf("Expected 1 broadcast for RequestChainMessage, got %d", len(mockPool.GetBroadcasts()))
	}

	broadcast := mockPool.GetBroadcasts()[0]
	msg, err := serializemessage.FromJSON([]byte(broadcast))
	if err != nil {
		t.Fatalf("Failed to deserialize broadcast message: %v", err)
	}
	if msg.Content.MessageType() != message.RequestLastNBlocksMessage.String() {
		t.Fatalf("Expected RequestMessageLastNBlocks, got %v", msg.Content.MessageType())
	}
}

func TestP2PProtocol_RequestBlocksBeforeMessage(t *testing.T) {
	mockPool := NewMockConnectionPool()
	sender := make(chan message.Message, 1)
	protocol := p2pprotocol.New(mockPool, sender)

	protocol.RequestBlocksBeforeMessage()

	if len(mockPool.GetBroadcasts()) != 1 {
		t.Fatalf("Expected 1 broadcast for RequestBlocksBeforeMessage, got %d", len(mockPool.GetBroadcasts()))
	}

	broadcast := mockPool.GetBroadcasts()[0]
	msg, err := serializemessage.FromJSON([]byte(broadcast))
	if err != nil {
		t.Fatalf("Failed to deserialize broadcast message: %v", err)
	}
	if msg.Content.MessageType() != message.RequestBlocksBeforeMessage.String() {
		t.Fatalf("Expected RequestMessageBlocksBefore, got %v", msg.Content.MessageType())
	}
}

func TestP2PProtocol_ResponseBlockMessage(t *testing.T) {
	mockPool := NewMockConnectionPool()
	sender := make(chan message.Message, 1)
	protocol := p2pprotocol.New(mockPool, sender)

	block := &block.Block{}
	protocol.ResponseBlockMessage(block)

	if len(mockPool.GetBroadcasts()) != 1 {
		t.Fatalf("Expected 1 broadcast for ResponseBlockMessage, got %d", len(mockPool.GetBroadcasts()))
	}

	broadcast := mockPool.GetBroadcasts()[0]
	msg, err := serializemessage.FromJSON([]byte(broadcast))
	if err != nil {
		t.Fatalf("Failed to deserialize broadcast message: %v", err)
	}
	if msg.Content.MessageType() != message.ResponseBlockMessage.String() {
		t.Fatalf("Expected ResponseMessageBlock, got %v", msg.Content.MessageType())
	}
}

func TestP2PProtocol_ResponseTransactionMessage(t *testing.T) {
	mockPool := NewMockConnectionPool()
	sender := make(chan message.Message, 1)
	protocol := p2pprotocol.New(mockPool, sender)

	tx := &transaction.Transaction{}
	protocol.ResponseTransactionMessage(tx)

	if len(mockPool.GetBroadcasts()) != 1 {
		t.Fatalf("Expected 1 broadcast for ResponseTransactionMessage, got %d", len(mockPool.GetBroadcasts()))
	}

	broadcast := mockPool.GetBroadcasts()[0]
	msg, err := serializemessage.FromJSON([]byte(broadcast))
	if err != nil {
		t.Fatalf("Failed to deserialize broadcast message: %v", err)
	}
	if msg.Content.MessageType() != message.ResponseTransactionMessage.String() {
		t.Fatalf("Expected ResponseMessageTransaction, got %v", msg.Content.MessageType())
	}
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
