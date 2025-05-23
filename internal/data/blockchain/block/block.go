package block

import (
	"sender/internal/data/blockchain/transaction"
	"sender/internal/jsonutil"
)

type Block struct {
	ID           int                       `json:"id"`
	TimeCreated  int64                     `json:"time_create"`
	Transactions []transaction.Transaction `json:"transactions"`
	PreviousHash string                    `json:"previous_hash"`
	Nonce        uint64                    `json:"nonce"`
}

func (b *Block) ToJson() ([]byte, error) {
	if b.Transactions == nil {
		b.Transactions = []transaction.Transaction{}
	}
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
