package message

import "sender/internal/data/blockchain/transaction"

type TransactionMessage struct {
	BaseMessage
	Transaction *transaction.Transaction `json:"transaction"`
}
