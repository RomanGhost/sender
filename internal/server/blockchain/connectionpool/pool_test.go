package connectionpool

import (
	"net"
	"sender/internal/server/blockchain/connectionpool/message"
	"sender/internal/server/blockchain/connectionpool/peer"
	protocolmsg "sender/internal/server/blockchain/protocol/message"
	"strings"
	"sync"
	"testing"
	"time"
)

// setupPool returns a ConnectionPool with buffered channels for pool and protocol
func setupPool(timeoutSecs int64) (*ConnectionPool, chan message.PoolMessage, chan protocolmsg.Message) {
	poolChan := make(chan message.PoolMessage, 10)
	protocolChan := make(chan protocolmsg.Message, 10)
	cp := NewConnectionPool(poolChan, timeoutSecs, protocolChan)
	return &cp, poolChan, protocolChan
}

// TestAddRemovePeer verifies addConnection, getPeerAddresses, removeConnection
func TestAddRemovePeer(t *testing.T) {
	cp, _, _ := setupPool(10)
	addr := &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9000}
	c, s := net.Pipe()
	defer c.Close()
	defer s.Close()
	prot := &peer.ProtectedConnection{Conn: s, Mutex: &sync.Mutex{}}
	cp.addConnection(addr, prot)

	addrs := cp.getPeerAddresses()
	if len(addrs) != 1 || addrs[0].String() != addr.String() {
		t.Fatalf("Expected peer %s, got %v", addr, addrs)
	}

	cp.removeConnection(addr)
	if len(cp.getPeerAddresses()) != 0 {
		t.Fatalf("Expected no peers after removal, got %v", cp.getPeerAddresses())
	}
}

// TestSendToPeer verifies sending a message to a single peer without blocking
func TestSendToPeer(t *testing.T) {
	cp, _, _ := setupPool(10)
	addr := &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9001}
	client, server := net.Pipe()
	defer client.Close()
	defer server.Close()

	prot := &peer.ProtectedConnection{Conn: server, Mutex: &sync.Mutex{}}
	cp.addConnection(addr, prot)
	msg := "hello"

	// concurrent read
	lineCh := make(chan string)
	go func() {
		buf := make([]byte, 64)
		n, _ := client.Read(buf)
		lineCh <- strings.TrimSpace(string(buf[:n]))
	}()

	if err := cp.sendToPeer(addr, msg); err != nil {
		t.Fatalf("sendToPeer failed: %v", err)
	}

	line := <-lineCh
	if line != msg {
		t.Errorf("Expected '%s', got '%s'", msg, line)
	}
}

// TestBroadcastAndRemoveFailed ensures broadcast writes to good peers and removes failed ones
func TestBroadcastAndRemoveFailed(t *testing.T) {
	cp, _, _ := setupPool(10)

	// good peer
	addr1 := &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9002}
	c1, s1 := net.Pipe()
	defer c1.Close()
	defer s1.Close()
	cp.addConnection(addr1, &peer.ProtectedConnection{Conn: s1, Mutex: &sync.Mutex{}})

	// bad peer (client end closed)
	addr2 := &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9003}
	c2, s2 := net.Pipe()
	c2.Close()
	defer s2.Close()
	cp.addConnection(addr2, &peer.ProtectedConnection{Conn: s2, Mutex: &sync.Mutex{}})

	msg := "broadcast"
	// concurrent read for good peer
	lineCh := make(chan string)
	go func() {
		buf := make([]byte, 64)
		n, _ := c1.Read(buf)
		lineCh <- strings.TrimSpace(string(buf[:n]))
	}()

	cp.broadcast(msg)

	line := <-lineCh
	if line != msg {
		t.Errorf("Good peer expected '%s', got '%s'", msg, line)
	}

	addrs := cp.getPeerAddresses()
	if len(addrs) != 1 || addrs[0].String() != addr1.String() {
		t.Errorf("Expected only %s, got %v", addr1, addrs)
	}
}

// TestCleanupInactive prunes peers not seen within timeout
func TestCleanupInactive(t *testing.T) {
	cp, _, _ := setupPool(0)
	addr := &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9004}
	_, s := net.Pipe()
	defer s.Close()
	cp.addConnection(addr, &peer.ProtectedConnection{Conn: s, Mutex: &sync.Mutex{}})

	// simulate outdated LastSeen
	cp.mutex.Lock()
	cp.connections[addr.String()].LastSeen = time.Now().Add(-2 * time.Second)
	cp.mutex.Unlock()

	cp.cleanupInactive()
	if len(cp.getPeerAddresses()) != 0 {
		t.Errorf("Expected no peers after cleanup, got %v", cp.getPeerAddresses())
	}
}

// TestHandlePeerMessage splits data on newline and forwards with RawMessage
func TestHandlePeerMessage(t *testing.T) {
	cp, _, proto := setupPool(10)
	addr := &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9005}
	_, s := net.Pipe()
	defer s.Close()
	cp.addConnection(addr, &peer.ProtectedConnection{Conn: s, Mutex: &sync.Mutex{}})

	// chunks
	cp.handlePeerMessage(addr, "one\ntwo\npar")
	cp.handlePeerMessage(addr, "tial\n")

	var got []string
	for i := 0; i < 3; i++ {
		select {
		case pm := <-proto:
			rm := pm.Content.(*protocolmsg.RawMessage)
			got = append(got, string(rm.MessageJson))
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Expected RawMessage")
		}
	}

	exp := []string{"one", "two", "partial"}
	for i := range exp {
		if got[i] != exp[i] {
			t.Errorf("Expected %s, got %s", exp[i], got[i])
		}
	}
}

// TestGetPeers verifies GetPeers via poolChan and response channel
func TestGetPeers(t *testing.T) {
	cp, poolChan, _ := setupPool(10)
	addr1 := &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9006}
	addr2 := &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9007}
	_, s1 := net.Pipe()
	_, s2 := net.Pipe()
	defer s1.Close()
	defer s2.Close()
	cp.addConnection(addr1, &peer.ProtectedConnection{Conn: s1, Mutex: &sync.Mutex{}})
	cp.addConnection(addr2, &peer.ProtectedConnection{Conn: s2, Mutex: &sync.Mutex{}})

	resp := make(chan []net.Addr, 1)
	poolChan <- message.PoolMessage{Type: message.GetPeers, ResponseChan: resp}
	go cp.Run()

	addrs := <-resp
	if len(addrs) != 2 {
		t.Errorf("Expected 2 peers, got %d", len(addrs))
	}
}

// TestRunNewPeerAndProtocol ensures Run handles NewPeer messages
func TestRunNewPeerAndProtocol(t *testing.T) {
	cp, poolChan, proto := setupPool(10)
	addr := &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9008}
	_, s := net.Pipe()
	defer s.Close()
	poolChan <- message.PoolMessage{Type: message.NewPeer, Addr: addr, Conn: &peer.ProtectedConnection{Conn: s, Mutex: &sync.Mutex{}}}

	go cp.Run()

	first := <-proto
	second := <-proto
	if first.Type != protocolmsg.ResponsePeerMessage {
		t.Errorf("Expected ResponsePeerMessage, got %v", first.Type)
	}
	if second.Type != protocolmsg.ResponseMessageInfo {
		t.Errorf("Expected ResponseMessageInfo, got %v", second.Type)
	}
}
