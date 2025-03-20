package message

import (
	"net"
	"sender/internal/server/blockchain/connectionpool/peer"
)

// MessageType defines the type of pool message
type MessageType int

const (
	NewPeer MessageType = iota
	PeerDisconnected
	BroadcastMessage
	GetPeers
	PeerMessage
)

// PoolMessage represents a message to the connection pool
type PoolMessage struct {
	Type         MessageType
	Addr         net.Addr
	Conn         *peer.ProtectedConnection
	Message      string
	ResponseChan chan []net.Addr // For GetPeers responses
}
