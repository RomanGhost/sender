package server

import (
	"fmt"
	"log"
	"net"
	"sender/internal/server/blockchain/connectionPool/message"
	"sender/internal/server/blockchain/connectionpool/peer"
	"sync"
	"time"
)

// Server represents the P2P server that listens for incoming connections
type Server struct {
	poolChan chan message.PoolMessage
}

// NewServer creates a new P2P server instance
func NewServer(poolChan chan message.PoolMessage) *Server {
	return &Server{
		poolChan: poolChan,
	}
}

// Run starts the server and begins listening for connections
func (s *Server) Run(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Printf("P2P server started on %s", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		localAddr := conn.LocalAddr().(*net.TCPAddr)
		remoteAddr := conn.RemoteAddr().(*net.TCPAddr)

		if localAddr.IP.Equal(remoteAddr.IP) && localAddr.Port == remoteAddr.Port {
			conn.Close()
			continue
		}

		// Start a new goroutine for each peer
		go s.handle(conn)
	}
}

// Connect attempts to connect to a peer at the given address
func (s *Server) Connect(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Printf("Error connecting to %s: %v", address, err)
		return err
	}

	localAddr := conn.LocalAddr().(*net.TCPAddr)
	remoteAddr := conn.RemoteAddr().(*net.TCPAddr)

	if localAddr.IP.Equal(remoteAddr.IP) && localAddr.Port == remoteAddr.Port {
		conn.Close()
		return fmt.Errorf("attempted to connect to self")
	}

	// Start a new goroutine for the connection
	go s.handle(conn)
	return nil
}

// GetPoolSender returns a channel for sending messages to the connection pool
func (s *Server) GetPoolSender() chan<- message.PoolMessage {
	return s.poolChan
}

// handle manages a single connection
func (s *Server) handle(conn net.Conn) error {
	addr := conn.RemoteAddr()
	log.Printf("Started thread for peer %s", addr.String())

	// Create a mutex-protected connection
	connMutex := &sync.Mutex{}
	wrappedConn := peer.NewProtectedConnection(conn, connMutex)

	// Notify the pool about the new peer
	s.poolChan <- message.PoolMessage{
		Type: message.NewPeer,
		Addr: addr,
		Conn: wrappedConn,
	}

	// Set read timeout
	conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))

	buffer := make([]byte, 1024)
	for {
		// Read data from the peer
		connMutex.Lock()
		n, err := conn.Read(buffer)
		connMutex.Unlock()

		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// Timeout occurred, reset the deadline and continue
				conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
				time.Sleep(100 * time.Millisecond)
				continue
			} else if err.Error() == "EOF" {
				// Connection closed
				log.Printf("Peer %s disconnected", addr.String())
				break
			} else {
				// Other error
				log.Printf("Error reading from peer %s: %v", addr.String(), err)
				break
			}
		}

		if n > 0 {
			message_json := string(buffer[:n])
			// Send the message to the pool
			s.poolChan <- message.PoolMessage{
				Type:    message.PeerMessage,
				Addr:    addr,
				Message: message_json,
			}
		}
	}

	// Notify the pool about the disconnected peer
	s.poolChan <- message.PoolMessage{
		Type: message.PeerDisconnected,
		Addr: addr,
	}

	return nil
}
