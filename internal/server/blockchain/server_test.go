package blockchain_test

import (
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"sender/internal/server/blockchain"
	"sender/internal/server/blockchain/connectionpool/message"
)

func startServer(t *testing.T, poolChan chan message.PoolMessage, addr string) {
	s := blockchain.NewServer(poolChan)
	go func() {
		err := s.Run(addr)
		if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			t.Errorf("Server Run failed: %v", err)
		}
	}()
	// Подождем немного, чтобы сервер запустился
	time.Sleep(100 * time.Millisecond)
}

func TestConnectAndSend(t *testing.T) {
	poolChan := make(chan message.PoolMessage, 10)
	addr := "127.0.0.1:9010"

	// Запускаем сервер
	startServer(t, poolChan, addr)

	s := blockchain.NewServer(poolChan)

	// Подключаемся к серверу
	err := s.Connect(addr)
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	// Ждем появления NewPeer
	var msg message.PoolMessage
	select {
	case msg = <-poolChan:
	case <-time.After(time.Second):
		t.Fatal("Timeout waiting for NewPeer message")
	}

	if msg.Type != message.NewPeer {
		t.Errorf("Expected NewPeer, got %v", msg.Type)
	}
}

func TestHandlePeerMessageAndDisconnect(t *testing.T) {
	poolChan := make(chan message.PoolMessage, 10)
	addr := "127.0.0.1:9011"

	startServer(t, poolChan, addr)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("Dial failed: %v", err)
	}

	// Ждем сообщения NewPeer
	select {
	case msg := <-poolChan:
		if msg.Type != message.NewPeer {
			t.Fatalf("Expected NewPeer, got %v", msg.Type)
		}
	case <-time.After(time.Second):
		t.Fatal("No NewPeer received")
	}

	// Отправляем данные
	_, err = conn.Write([]byte("test message"))
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	// Проверяем PeerMessage
	select {
	case msg := <-poolChan:
		if msg.Type != message.PeerMessage {
			t.Fatalf("Expected PeerMessage, got %v", msg.Type)
		}
		if msg.Message != "test message" {
			t.Errorf("Expected 'test message', got '%s'", msg.Message)
		}
	case <-time.After(time.Second):
		t.Fatal("No PeerMessage received")
	}

	// Закрываем соединение
	conn.Close()

	// Ждем сообщения PeerDisconnected
	select {
	case msg := <-poolChan:
		if msg.Type != message.PeerDisconnected {
			t.Fatalf("Expected PeerDisconnected, got %v", msg.Type)
		}
	case <-time.After(time.Second):
		t.Fatal("No PeerDisconnected received")
	}
}

func TestGetPoolSender(t *testing.T) {
	poolChan := make(chan message.PoolMessage, 1)
	s := blockchain.NewServer(poolChan)

	sender := s.GetPoolSender()

	go func() {
		sender <- message.PoolMessage{Type: message.GetPeers}
	}()

	select {
	case msg := <-poolChan:
		if msg.Type != message.GetPeers {
			t.Errorf("Expected GetPeers message, got %v", msg.Type)
		}
	case <-time.After(time.Second):
		t.Fatal("Did not receive message from GetPoolSender")
	}
}

func TestConnectToSelf(t *testing.T) {
	poolChan := make(chan message.PoolMessage, 1)
	addr := "127.0.0.1:9012"

	startServer(t, poolChan, addr)

	s := blockchain.NewServer(poolChan)

	// Открываем соединение напрямую, чтобы узнать локальный адрес
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("Dial failed: %v", err)
	}
	defer conn.Close()

	// Получаем IP и порт, на который мы подключены
	localAddr := conn.LocalAddr().(*net.TCPAddr)

	// Теперь пробуем подключиться к себе (тот же IP и порт)
	selfAddr := net.JoinHostPort(localAddr.IP.String(), fmt.Sprintf("%d", localAddr.Port))
	err = s.Connect(selfAddr)
	if err == nil {
		t.Fatal("Expected error when connecting to self, got nil")
	}
	if !strings.Contains(err.Error(), "No connection could be made because the target machine actively refused it") {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestConnectToInvalidAddress(t *testing.T) {
	poolChan := make(chan message.PoolMessage, 1)
	s := blockchain.NewServer(poolChan)

	err := s.Connect("127.0.0.1:9999") // Этот порт не слушает
	if err == nil {
		t.Fatal("Expected error when connecting to invalid address, got nil")
	}
}
