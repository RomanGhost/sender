package message

import "sender/internal/data/blockchain/block"

type ChainMessage struct {
	BaseMessage
	Chain []block.Block `json:"chain"`
}
