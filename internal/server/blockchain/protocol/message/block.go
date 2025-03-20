package message

import "sender/internal/data/blockchain/block"

type BlockMessage struct {
	BaseMessage
	Block *block.Block `json:"block"`
}
