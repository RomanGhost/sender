package responce

import (
	"sender/data/blockchain/block"
	"sender/internal/server/p2pprotocol/message"
	"time"
)

type BlockMessage struct {
	message.BaseMessage
	Block *block.Block `json:"block"`
	Force bool         `json:"force"`
}

func NewBlockMessage(block *block.Block, force bool) *BlockMessage {
	return &BlockMessage{
		BaseMessage: message.BaseMessage{
			ID:        0,
			Timestamp: time.Now(),
		},
		Block: block,
		Force: force,
	}
}

func (m *BlockMessage) MessageType() string {
	return message.ResponseBlockMessage.String()
}
