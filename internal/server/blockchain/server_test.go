package blockchain

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"testing"
	"time"

	"sender/internal/server/blockchain/p2pprotocol/message"
)

func TestNewBlockchainServer(t *testing.T) {
	sender := make(chan message.Message, 1)
	server := New("127.0.0.1", 8080, sender)

	if server.Address != "127.0.0.1" {
		t.Errorf("Expected Address to be '127.0.0.1', got '%s'", server.Address)
	}
	if server.Port != 8080 {
		t.Errorf("Expected Port to be 8080, got %d", server.Port)
	}
	if server.ConnectionPool == nil {
		t.Error("ConnectionPool should not be nil")
	}
	if server.P2PProtocol == nil {
		t.Error("P2PProtocol should not be nil")
	}
}

func TestRun_AcceptsConnection(t *testing.T) {
	sender := make(chan message.Message, 1)
	server := New("127.0.0.1", 8080, sender)

	go func() {
		server.Run()
	}()

	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	handshake, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("Failed to read handshake: %v", err)
	}
	if handshake != HandshakeMessage {
		t.Errorf("Expected handshake '%s', but got '%s'", HandshakeMessage, handshake)
	}
}

func TestHandleConnection_Timeout(t *testing.T) {
	t.Skip("Skipping this test temporarily")

	sender := make(chan message.Message, 1)
	main_server := New("127.0.0.1", 19090, sender)
	go func() {
		main_server.Run()
	}()

	client_server := New("127.0.0.1", 19091, sender)
	go func() {
		client_server.Run()
	}()

	err := client_server.Connect("127.0.0.1", 19090)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}

	<-time.After(Timeout * time.Minute)

	time.Sleep(3 * time.Second)
	if len(main_server.ConnectionPool.GetAlivePeers()) != 0 {
		t.Errorf("Expected 0 peers after timeout, got %d", len(main_server.ConnectionPool.GetAlivePeers()))
	}

	if len(client_server.ConnectionPool.GetAlivePeers()) != 0 {
		t.Errorf("Expected 0 peers after timeout, got %d", len(client_server.ConnectionPool.GetAlivePeers()))
	}
}

func TestConnect(t *testing.T) {
	sender := make(chan message.Message, 1)
	main_server := New("127.0.0.1", 18080, sender)
	go func() {
		main_server.Run()
	}()

	client_server := New("127.0.0.1", 18081, sender)
	go func() {
		client_server.Run()
	}()

	err := client_server.Connect("127.0.0.1", 18080)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	time.Sleep(time.Second)
	fmt.Println(main_server.ConnectionPool.GetAlivePeers())

	if len(main_server.ConnectionPool.GetAlivePeers()) != 1 {
		t.Errorf("Expected 1 peer in connection pool, got %d", len(main_server.ConnectionPool.GetAlivePeers()))
	}
}

func TestHandleConnection(t *testing.T) {
	sender := make(chan message.Message, 1)
	server := New("127.0.0.1", 8080, sender)

	mockConn := NewMockConn()
	mockConn.Reader = bufio.NewReader(bytes.NewReader([]byte(HandshakeMessage)))
	mockConn.Writer = new(bytes.Buffer)

	go server.handleConnection(mockConn)

	<-time.After(100 * time.Millisecond)

	if mockConn.Writer.String() != HandshakeMessage {
		t.Errorf("Expected handshake '%s', got '%s'", HandshakeMessage, mockConn.Writer.String())
	}

	if len(server.ConnectionPool.GetAlivePeers()) != 1 {
		t.Errorf("Expected 1 peer in connection pool, got %d", len(server.ConnectionPool.GetAlivePeers()))
	}
}

func TestExtractMessage(t *testing.T) {
	data := "message1\nmessage2\n"
	message, remaining := extractMessage(data)

	if message != "message1" {
		t.Errorf("Expected 'message1', got '%s'", message)
	}
	if remaining != "message2\n" {
		t.Errorf("Expected remaining to be 'message2\\n', got '%s'", remaining)
	}

	// Case: No newline in data
	data = "incompleteMessage"
	message, remaining = extractMessage(data)
	if message != "" {
		t.Errorf("Expected empty message, got '%s'", message)
	}
	if remaining != "incompleteMessage" {
		t.Errorf("Expected remaining to be 'incompleteMessage', got '%s'", remaining)
	}
}

func TestFindNewlineIndex(t *testing.T) {
	data := "hello\nworld"
	index := findNewlineIndex(data)
	if index != 5 {
		t.Errorf("Expected index 5, got %d", index)
	}

	data = "nowordnewline"
	index = findNewlineIndex(data)
	if index != -1 {
		t.Errorf("Expected -1, got %d", index)
	}
}
