package connectionpool

import (
	"fmt"
	"net"
	"sync"
)

type ConnectionPool struct {
	mu         sync.Mutex
	peers      map[string]net.Conn
	bufferSize int
}

func New(bufferSize int) *ConnectionPool {
	return &ConnectionPool{
		peers:      make(map[string]net.Conn),
		bufferSize: bufferSize,
	}
}

// AddPeer добавляет нового пира в пул.
func (cp *ConnectionPool) AddPeer(address string, conn net.Conn) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	cp.peers[address] = conn
	fmt.Printf("Added peer: %s\n", address)
}

// RemovePeer удаляет пира из пула.
func (cp *ConnectionPool) RemovePeer(address string) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	delete(cp.peers, address)
	fmt.Printf("Removed peer: %s\n", address)
}

// GetAlivePeers возвращает список активных соединений.
func (cp *ConnectionPool) GetAlivePeers() []net.Conn {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	peers := make([]net.Conn, 0, len(cp.peers))
	for _, conn := range cp.peers {
		peers = append(peers, conn)
	}
	return peers
}

// GetPeerAddresses возвращает список адресов подключенных пиров.
func (cp *ConnectionPool) GetPeerAddresses() []string {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	addresses := make([]string, 0, len(cp.peers))
	for addr := range cp.peers {
		addresses = append(addresses, addr)
	}
	return addresses
}

// Broadcast отправляет сообщение всем подключенным пирами.
func (cp *ConnectionPool) Broadcast(message string) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	disconnectedPeers := []string{}
	bufferSize := cp.bufferSize

	// Добавляем перенос строки для разделения сообщений.
	message += "\n"
	startIndex := 0

	// Разбиваем сообщение на части в зависимости от размера буфера.
	for startIndex < len(message) {
		endIndex := startIndex + bufferSize
		if endIndex > len(message) {
			endIndex = len(message)
		}
		messageChunk := message[startIndex:endIndex]
		// fmt.Printf("Само сообщение: %s\n", messageChunk)

		// Отправляем сообщение каждому пиру.
		for address, conn := range cp.peers {
			if _, err := conn.Write([]byte(messageChunk)); err != nil {
				fmt.Printf("Failed to send message to %s: %v\n", address, err)
				disconnectedPeers = append(disconnectedPeers, address)
			}
		}
		startIndex += bufferSize
	}

	// Удаляем отключившихся пиров.
	for _, address := range disconnectedPeers {
		cp.RemovePeer(address)
	}
}

// GetBuffer создает пустой буфер заданного размера.
func (cp *ConnectionPool) GetBuffer() []byte {
	return make([]byte, cp.bufferSize)
}
