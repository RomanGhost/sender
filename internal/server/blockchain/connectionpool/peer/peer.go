package peer

import (
	"net"
	"sync"
	"time"
)

// ProtectedConnection is a wrapper around a connection with mutex protection
type ProtectedConnection struct {
	Conn  net.Conn
	Mutex *sync.Mutex
}

// NewProtectedConnection creates a new protected connection
func NewProtectedConnection(conn net.Conn, mutex *sync.Mutex) ProtectedConnection {
	return ProtectedConnection{
		Conn:  conn,
		Mutex: mutex,
	}
}

// PeerConnection represents a connection to a peer
type PeerConnection struct {
	Addr     net.Addr
	Conn     *ProtectedConnection
	LastSeen time.Time
	Buffer   string
}
