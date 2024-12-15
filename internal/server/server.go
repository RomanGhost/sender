package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"

	"sender/internal/p2pprotocol"
	"sender/internal/p2pprotocol/message"
	"sender/internal/server/connectionpool"
)

const HandshakeMessage = "NEW_CONNECT!\r\n"
const Timeout = 5 //mins
const BufferSize = 4096

type BlockchainServer struct {
	Address        string
	Port           int
	ConnectionPool *connectionpool.ConnectionPool
	P2PProtocol    *p2pprotocol.P2PProtocol
}

func New(address string, port int, sender chan message.Message) *BlockchainServer {
	connectionPool := connectionpool.New(BufferSize)
	p2pProtocol := p2pprotocol.New(connectionPool, sender)

	return &BlockchainServer{
		Address:        address,
		Port:           port,
		ConnectionPool: connectionPool,
		P2PProtocol:    p2pProtocol,
	}
}

func (bs *BlockchainServer) Run() {
	listenerAddress := fmt.Sprintf("%s:%d", bs.Address, bs.Port)
	listener, err := net.Listen("tcp", listenerAddress)
	if err != nil {
		log.Fatalf("Error: %v", err)
		return
	}
	defer listener.Close()

	log.Printf("Server started on %s", listenerAddress)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go bs.handleConnection(conn)
	}
}

func (bs *BlockchainServer) Connect(address string, port int) {
	connectionAddress := fmt.Sprintf("%s:%d", address, port)
	conn, err := net.Dial("tcp", connectionAddress)
	if err != nil {
		log.Printf("Error connecting to server: %v", err)
		return
	}

	log.Printf("Connected to %s", connectionAddress)

	go bs.handleConnection(conn)
}

func (bs *BlockchainServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	peerAddress := conn.RemoteAddr().String()
	log.Printf("New connection from %s", peerAddress)

	reader := bufio.NewReader(conn)
	conn.Write([]byte(HandshakeMessage))

	// Check handshake
	handshake, err := reader.ReadString('\n')
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return // Время ожидания истекло
		} else if err.Error() == "use of closed network connection" {
			log.Printf("Connection closed by peer: %s", peerAddress)
			bs.ConnectionPool.RemovePeer(peerAddress)
			return
		}
		log.Printf("Failed to read handshake from %s: %v", peerAddress, err)
		return
	}

	if handshake != HandshakeMessage {
		log.Printf("Unauthorized client from %s, message: %v == %v", peerAddress, []byte(handshake), []byte(HandshakeMessage))
		conn.Write([]byte("Unauthorized\n"))
		return
	}

	log.Printf("Authorized client connected from %s", peerAddress)
	bs.ConnectionPool.AddPeer(peerAddress, conn)

	// Initialize last message time
	lastMessageTime := time.Now()

	buffer := make([]byte, BufferSize)
	var accumulatedData string

	for {
		// Check for timeout
		if time.Since(lastMessageTime) > Timeout*time.Minute {
			log.Printf("Client %s inactive for %v minutes, disconnecting", peerAddress, Timeout)
			bs.ConnectionPool.RemovePeer(peerAddress)
			break
		}

		// Set read deadline for inactivity check
		conn.SetReadDeadline(time.Now().Add(Timeout * time.Second))

		n, err := conn.Read(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue // Ignore timeouts
			}
			log.Printf("Error reading from %s: %v", peerAddress, err)
			bs.ConnectionPool.RemovePeer(peerAddress)
			break
		}

		lastMessageTime = time.Now()
		accumulatedData += string(buffer[:n])

		// Process messages
		for {
			message, remainingData := extractMessage(accumulatedData)
			if message == "" {
				break
			}
			accumulatedData = remainingData
			log.Printf("Received message from %s: %s", peerAddress, message)
			// Processing logic for the message would go here
		}
	}
}

// Extracts a message from a string of data, separated by '\n'
func extractMessage(data string) (string, string) {
	if index := findNewlineIndex(data); index != -1 {
		return data[:index], data[index+1:]
	}
	return "", data
}

// Finds the index of the newline character in the string
func findNewlineIndex(data string) int {
	for i, c := range data {
		if c == '\n' {
			return i
		}
	}
	return -1
}
