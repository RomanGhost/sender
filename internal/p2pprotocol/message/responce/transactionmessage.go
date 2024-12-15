package responce

import (
	"sender/internal/data/transaction"
	"sender/internal/p2pprotocol/message"
	"time"
)

type TransactionMessage struct {
	message.BaseMessage
	Transaction *transaction.Transaction `json:"transaction"` // Упрощённо: строка вместо структуры транзакции
}

func NewTransactionMessage(transaction *transaction.Transaction) *TransactionMessage {
	return &TransactionMessage{
		BaseMessage: message.BaseMessage{
			ID:        0,
			Timestamp: time.Now(),
		},
		Transaction: transaction,
	}
}

func (m *TransactionMessage) MessageType() string {
	return "ResponseTransactionMessage"
}
