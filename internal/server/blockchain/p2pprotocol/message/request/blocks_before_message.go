package request

import (
	"sender/internal/server/blockchain/p2pprotocol/message"
	"time"
)

type BlocksBeforeMessage struct {
	message.BaseMessage
}

func NewBlocksBeforeMessage() *BlocksBeforeMessage {
	return &BlocksBeforeMessage{
		BaseMessage: message.BaseMessage{
			ID:        0,
			Timestamp: time.Now(),
		},
	}
}

func (fm *BlocksBeforeMessage) MessageType() string {
	return message.RequestBlocksBeforeMessage.String()
}
