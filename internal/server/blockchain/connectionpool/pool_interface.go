package connectionpool

import "net"

type ConnectionPoolInterface interface {
	AddPeer(address string, conn net.Conn)
	RemovePeer(address string)
	GetPeerAddresses() []string
	Broadcast(message string)
}
