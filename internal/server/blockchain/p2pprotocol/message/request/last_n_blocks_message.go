package request

import (
	"sender/internal/server/blockchain/p2pprotocol/message"
	"time"
)

type LastNBlocksMessage struct {
	message.BaseMessage
	N uint `json:"n"`
}

func NewLastNBlocksMessage(countBlocks uint) *LastNBlocksMessage {
	return &LastNBlocksMessage{
		BaseMessage: message.BaseMessage{
			ID:        0,
			Timestamp: time.Now(),
		},
		N: countBlocks,
	}
}

func (m *LastNBlocksMessage) MessageType() string {
	return message.RequestLastNBlocksMessage.String()
}
