package connectionpool

import (
	"fmt"
	"log"
	"net"
	"sender/internal/server/blockchain/connectionpool/message"
	"sender/internal/server/blockchain/connectionpool/peer"
	"sender/internal/server/blockchain/protocol"
	"strings"
	"sync"
	"time"
)

// ConnectionPool manages all peer connections
type ConnectionPool struct {
	connections map[string]*peer.PeerConnection
	mutex       sync.RWMutex
	timeout     time.Duration

	// Channels for pool communication
	poolChan chan message.PoolMessage

	// Channel for communication with the protocol
	protocolChan chan<- protocol.Message
}

// NewConnectionPool creates a new connection pool
func NewConnectionPool(timeoutSecs int64, protocolChan chan<- protocol.Message) *ConnectionPool {
	return &ConnectionPool{
		connections:  make(map[string]*peer.PeerConnection),
		timeout:      time.Duration(timeoutSecs) * time.Second,
		poolChan:     make(chan message.PoolMessage, 100),
		protocolChan: protocolChan,
	}
}

// GetPoolChan returns the channel for sending messages to the pool
func (cp *ConnectionPool) GetPoolChan() chan<- message.PoolMessage {
	return cp.poolChan
}

// addConnection adds a new peer connection to the pool
func (cp *ConnectionPool) addConnection(addr net.Addr, conn *peer.ProtectedConnection) {
	addrStr := addr.String()

	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	cp.connections[addrStr] = &peer.PeerConnection{
		Addr:     addr,
		Conn:     conn,
		LastSeen: time.Now(),
		Buffer:   "",
	}

	log.Printf("New peer connected: %s, total peers: %d", addrStr, len(cp.connections))
}

// removeConnection removes a peer connection from the pool
func (cp *ConnectionPool) removeConnection(addr net.Addr) {
	addrStr := addr.String()

	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	if _, exists := cp.connections[addrStr]; exists {
		delete(cp.connections, addrStr)
		log.Printf("Peer removed: %s", addrStr)
	}
}

// getPeerAddresses returns a list of all peer addresses
func (cp *ConnectionPool) getPeerAddresses() []net.Addr {
	cp.mutex.RLock()
	defer cp.mutex.RUnlock()

	addrs := make([]net.Addr, 0, len(cp.connections))
	for _, peer := range cp.connections {
		addrs = append(addrs, peer.Addr)
	}

	return addrs
}

// sendToPeer sends a message to a specific peer
func (cp *ConnectionPool) sendToPeer(addr net.Addr, message string) error {
	addrStr := addr.String()

	cp.mutex.RLock()
	peer, exists := cp.connections[addrStr]
	cp.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("peer not found: %s", addrStr)
	}

	peer.Conn.Mutex.Lock()
	defer peer.Conn.Mutex.Unlock()

	if _, err := fmt.Fprintf(peer.Conn.Conn, "%s\n", message); err != nil {
		return err
	}

	peer.LastSeen = time.Now()
	return nil
}

// broadcast sends a message to all peers
func (cp *ConnectionPool) broadcast(message string) {
	cp.mutex.RLock()
	peers := make([]*peer.PeerConnection, 0, len(cp.connections))
	for _, peer := range cp.connections {
		peers = append(peers, peer)
	}
	cp.mutex.RUnlock()

	var failedPeers []net.Addr

	for _, peer := range peers {
		peer.Conn.Mutex.Lock()
		_, err := fmt.Fprintf(peer.Conn.Conn, "%s\n", message)
		peer.Conn.Mutex.Unlock()

		if err != nil {
			failedPeers = append(failedPeers, peer.Addr)
		} else {
			peer.LastSeen = time.Now()
		}
	}

	// Remove failed peers
	for _, addr := range failedPeers {
		cp.removeConnection(addr)
	}
}

// cleanupInactive removes inactive connections
func (cp *ConnectionPool) cleanupInactive() {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	now := time.Now()
	var inactivePeers []string

	for addrStr, peer := range cp.connections {
		if now.Sub(peer.LastSeen) > cp.timeout {
			inactivePeers = append(inactivePeers, addrStr)
		}
	}

	for _, addrStr := range inactivePeers {
		log.Printf("Inactive peer timeout: %s", addrStr)
		delete(cp.connections, addrStr)
	}
}

// Run starts the connection pool message processing
func (cp *ConnectionPool) Run() {
	for {
		select {
		case msg := <-cp.poolChan:
			switch msg.Type {
			case message.NewPeer:
				cp.addConnection(msg.Addr, msg.Conn)

				// Notify the protocol about the new peer
				peerMsg := protocol.NewPeerMessage(msg.Addr.(*net.TCPAddr).IP.String())
				cp.protocolChan <- peerMsg

				// Request initial message info
				cp.protocolChan <- protocol.NewInfoMessage()

			case message.PeerDisconnected:
				cp.removeConnection(msg.Addr)

			case message.BroadcastMessage:
				log.Printf("Broadcasting message: %s", msg.Message)
				cp.broadcast(msg.Message)

			case message.GetPeers:
				peers := cp.getPeerAddresses()
				msg.ResponseChan <- peers

			case message.PeerMessage:
				cp.handlePeerMessage(msg.Addr, msg.Message)
			}

		case <-time.After(600 * time.Second):
			log.Printf("Pool cleanup triggered")
			cp.cleanupInactive()
		}
	}
}

// handlePeerMessage processes messages from peers
func (cp *ConnectionPool) handlePeerMessage(addr net.Addr, message string) {
	addrStr := addr.String()

	cp.mutex.Lock()
	peer, exists := cp.connections[addrStr]
	if !exists {
		cp.mutex.Unlock()
		log.Printf("Message from unknown peer: %s", addrStr)
		return
	}

	// Append to buffer and process
	buffer := peer.Buffer + message
	peer.Buffer = "" // Clear the buffer
	cp.mutex.Unlock()

	messages := []string{}

	// Split the buffer into messages by newline
	for {
		parts := strings.SplitN(buffer, "\n", 2)
		if len(parts) == 1 {
			// No more newlines, store the rest in the buffer
			cp.mutex.Lock()
			peer.Buffer = parts[0]
			cp.mutex.Unlock()
			break
		}

		// Add the message and continue processing
		messages = append(messages, parts[0])
		buffer = parts[1]
	}

	// Process the messages
	for _, msg := range messages {
		// Forward to the protocol
		cp.protocolChan <- protocol.NewRawMessage([]byte(msg))

		// Update last seen time
		cp.mutex.Lock()
		if peer, exists := cp.connections[addrStr]; exists {
			peer.LastSeen = time.Now()
		}
		cp.mutex.Unlock()
	}
}
