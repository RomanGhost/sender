package responce

import (
	"sender/internal/data/blockchain/block"
	"sender/internal/p2pprotocol/message"
	"time"
)

type ChainMessage struct {
	message.BaseMessage
	ChainMessage []block.Block
}

func NewChainMessage(chain []block.Block) *ChainMessage {
	return &ChainMessage{
		BaseMessage: message.BaseMessage{
			ID:        0,
			Timestamp: time.Now(),
		},
		ChainMessage: chain,
	}
}

func (cm *ChainMessage) MessageType() string {
	return "ResponseChainMessage"
}
