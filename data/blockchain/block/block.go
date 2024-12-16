package block

import (
	"sender/data/blockchain/transaction"
	"sender/internal/jsonutil"
	"time"
)

type Block struct {
	ID           int                       `json:"id"`
	TimeCreated  time.Time                 `json:"time_create"`
	Transactions []transaction.Transaction `json:"transactions"`
	PreviousHash string                    `json:"previous_hash"`
	Nonce        uint64                    `json:"nonce"`
}

func (b *Block) ToJson() ([]byte, error) {
	return jsonutil.ToJSON(b)
}

func FromJSON(blockJson []byte) (*Block, error) {
	var block Block
	err := jsonutil.FromJSON(blockJson, &block)

	return &block, err
}

func (b *Block) GetTransactions() []transaction.Transaction {
	return b.Transactions
}
