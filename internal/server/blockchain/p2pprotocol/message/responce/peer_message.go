package responce

import (
	"sender/internal/server/blockchain/p2pprotocol/message"
	"time"
)

type PeerMessage struct {
	message.BaseMessage
	PeerAddresses []string `json:"peer_address"`
}

func NewPeerMessage(peer_addresses []string) *PeerMessage {
	return &PeerMessage{
		BaseMessage: message.BaseMessage{
			ID:        0,
			Timestamp: time.Now(),
		},
		PeerAddresses: peer_addresses,
	}
}

func (pm *PeerMessage) MessageType() string {
	return message.ResponsePeerMessage.String()
}
