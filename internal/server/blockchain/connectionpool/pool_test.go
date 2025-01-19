package connectionpool_test

import (
	"net"
	"sender/internal/server/blockchain/connectionpool"
	"sync"
	"testing"
	"time"
)

func TestConnectionPool_AddPeerAndGetPeerAddresses(t *testing.T) {
	pool := connectionpool.New(1024)

	// Создаем фейковое соединение
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	// Добавляем пира
	address := "peer1"
	pool.AddPeer(address, client)

	// Проверяем, что адрес появился в списке
	addresses := pool.GetPeerAddresses()
	if len(addresses) != 1 || addresses[0] != address {
		t.Errorf("Expected addresses to contain %s, got %v", address, addresses)
	}
}

func TestConnectionPool_RemovePeer(t *testing.T) {
	pool := connectionpool.New(1024)

	// Создаем фейковое соединение
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	// Добавляем пира
	address := "peer1"
	pool.AddPeer(address, client)

	// Удаляем пира
	pool.RemovePeer(address)

	// Проверяем, что список адресов пуст
	addresses := pool.GetPeerAddresses()
	if len(addresses) != 0 {
		t.Errorf("Expected addresses to be empty, got %v", addresses)
	}
}

func TestConnectionPool_Broadcast(t *testing.T) {
	pool := connectionpool.New(1024)

	client1, server1 := net.Pipe()
	client2, server2 := net.Pipe()

	pool.AddPeer("peer1", client1)
	pool.AddPeer("peer2", client2)

	message := "Hello, peers!\n" // Добавляем перенос строки, как в Broadcast.

	var wg sync.WaitGroup

	// Функция для прослушивания сообщений от пиров.
	listen := func(server net.Conn, peerName string) {
		defer wg.Done()
		buf := make([]byte, 1024) // Создаем буфер достаточного размера.
		n, err := server.Read(buf)
		if err != nil {
			t.Errorf("Failed to read message from %s: %v", peerName, err)
			return
		}
		// Преобразуем буфер в строку.
		str := string(buf[:n])
		t.Logf("%s received: %s", peerName, str)

		// Проверяем, что сообщение соответствует ожидаемому.
		if str != message {
			t.Errorf("Expected message: %s, got: %s", message, str)
		}
	}

	wg.Add(2)
	go listen(server1, "peer1")
	go listen(server2, "peer2")

	// Вызываем Broadcast в главной горутине.
	time.Sleep(50 * time.Millisecond) // Ждем, чтобы серверы начали слушать.
	pool.Broadcast("Hello, peers!")   // Отправляем сообщение.

	wg.Wait()
}

func TestConnectionPool_BroadcastRemovePeer(t *testing.T) {
	pool := connectionpool.New(1024)

	client1, server1 := net.Pipe()
	client2, server2 := net.Pipe()

	pool.AddPeer("peer1", client1)
	pool.AddPeer("peer2", client2)
	client2.Close()

	message := "Hello, peers!\n" // Добавляем перенос строки, как в Broadcast.

	var wg sync.WaitGroup

	// Функция для прослушивания сообщений от пиров.
	accept := func(server net.Conn, peerName string) {
		defer wg.Done()
		buf := make([]byte, 1024) // Создаем буфер достаточного размера.
		n, err := server.Read(buf)
		if err != nil {
			t.Errorf("Failed to read message from %s: %v", peerName, err)
			return
		}
		// Преобразуем буфер в строку.
		str := string(buf[:n])
		t.Logf("%s received: %s", peerName, str)

		// Проверяем, что сообщение соответствует ожидаемому.
		if str != message {
			t.Errorf("Expected message: %s, got: %s", message, str)
		}
	}

	error_reason := func(server net.Conn) {
		defer wg.Done()
		buf := make([]byte, 1024) // Создаем буфер достаточного размера.
		_, err := server.Read(buf)
		if err == nil {
			t.Errorf("Failed to read message")
			return
		}
	}

	wg.Add(2)
	go accept(server1, "peer1")
	go error_reason(server2)

	// Вызываем Broadcast в главной горутине.
	time.Sleep(50 * time.Millisecond) // Ждем, чтобы серверы начали слушать.
	pool.Broadcast("Hello, peers!")   // Отправляем сообщение.

	wg.Wait()
}
