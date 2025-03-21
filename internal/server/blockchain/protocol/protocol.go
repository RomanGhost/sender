package protocol

import (
	"encoding/json"
	"log"
	"sender/internal/app"
	poolMessage "sender/internal/server/blockchain/connectionpool/message"
	"sender/internal/server/blockchain/protocol/message"
	"time"
)

// P2PProtocol manages the P2P communication protocol
type P2PProtocol struct {
	// Channels for protocol communication
	messageChan chan message.Message

	// Channel for communication with the connection pool
	poolChan chan<- poolMessage.PoolMessage

	lastMessageID uint64
	appState      *app.AppState
}

// NewP2PProtocol creates a new P2P protocol instance
func NewProtocol(messageChan chan message.Message, appState *app.AppState, poolChan chan<- poolMessage.PoolMessage) *P2PProtocol {
	return &P2PProtocol{
		messageChan:   messageChan, //make(chan message.Message, 100),
		poolChan:      poolChan,
		lastMessageID: 0,
		appState:      appState,
	}
}

// GetMessageChan returns the channel for sending messages to the protocol
func (p *P2PProtocol) GetMessageChan() chan<- message.Message {
	return p.messageChan
}

// Run starts the P2P protocol message processing
func (p *P2PProtocol) Run() {
	for {
		select {
		case msg := <-p.messageChan:
			switch msg.Type {
			case message.RawMessageType:
				rawMsgJson := msg.Content.(*message.RawMessage).MessageJson
				msg_from_json, err := message.MessageFromJson(rawMsgJson)
				if err != nil {
					log.Printf("Error with message")
					return
				}

				p.processMessage(*msg_from_json)

			default:
				// Message from this server
				p.sendMessage(msg)
			}

		case <-time.After(1 * time.Second):
			// Timeout, just continue
		}
	}
}

// processMessage handles incoming messages from peers
func (p *P2PProtocol) processMessage(msg message.Message) {
	// Check if we've already seen this message
	if msg.Type == message.RequestMessageInfo {
		log.Printf("Type:RequestMessageInfo received")
		p.sendFirstMessage()
		return
	}

	if msg.Type == message.ResponseMessageInfo {
		log.Printf("Type:ResponseMessageInfo received")

		if p.lastMessageID < msg.Content.GetID() {
			p.lastMessageID = msg.Content.GetID()
		}
		log.Printf("Received message info: %d/%d", msg.Content.GetID(), p.lastMessageID)
		return
	}

	// Check if this is a duplicate message
	if msg.Content.GetID() <= p.lastMessageID {
		log.Printf("Message ID less than current: %d<%d", msg.Content.GetID(), p.lastMessageID)
		return
	}

	// Update the last message ID
	p.lastMessageID = msg.Content.GetID()

	// Broadcast the message to all peers
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return
	}

	p.poolChan <- poolMessage.PoolMessage{
		Type:    poolMessage.BroadcastMessage,
		Message: string(msgJSON),
	}

	// Process the message based on its type
	switch msg.Type {
	case message.ResponseBlockMessage:
		blockMessage := msg.Content.(*message.BlockMessage)
		p.processBlock(blockMessage)

	case message.ResponsePeerMessage:
		peerMsg := msg.Content.(*message.PeerMessage)
		p.processPeer(peerMsg)

	case message.ResponseTextMessage:
		textMsg := msg.Content.(*message.TextMessage).Message
		log.Printf("Received text message: %s", textMsg)

	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}

// sendMessage sends a message to all peers
func (p *P2PProtocol) sendMessage(msg message.Message) {
	p.lastMessageID++
	msg.Content.SetID(p.lastMessageID)

	msgJSON, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal message: %v", err)
		return
	}

	p.poolChan <- poolMessage.PoolMessage{
		Type:    poolMessage.BroadcastMessage,
		Message: string(msgJSON),
	}
}

// processBlock processes a block message
func (p *P2PProtocol) processBlock(msg *message.BlockMessage) {
	p.appState.ReadBlock(msg)
}

// processPeer processes a peer message
func (p *P2PProtocol) processPeer(msg *message.PeerMessage) {
	peer := msg.PeerAddrIp
	log.Printf("Received new peer: %s", peer)

	p.appState.Connect(peer)
}

// sendFirstMessage sends the initial message after connecting to a peer
func (p *P2PProtocol) sendFirstMessage() {
	p.lastMessageID++
	responseMsg := message.NewInfoMessage()
	responseMsg.Content.SetID(p.lastMessageID)

	msgJSON, err := json.Marshal(responseMsg)
	if err != nil {
		log.Printf("Failed to marshal response message: %v", err)
		return
	}

	p.poolChan <- poolMessage.PoolMessage{
		Type:    poolMessage.BroadcastMessage,
		Message: string(msgJSON),
	}
}
