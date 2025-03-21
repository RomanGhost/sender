package message

import (
	"encoding/json"
	"log"
	"sender/internal/data/blockchain/block"
	"sender/internal/data/blockchain/transaction"
)

type MessageInterface interface {
	GetID() uint64
	SetID(newID uint64)
	GetTime() int64
}

// Message represents a P2P protocol message
type Message struct {
	Type    MessageType      `json:"type"`
	Content MessageInterface `json:"content"`
}

func MessageFromJson(messageJson []byte) (*Message, error) {
	var body struct {
		Type    MessageType     `json:"type"`
		Content json.RawMessage `json:"content"`
	}

	if err := json.Unmarshal(messageJson, &body); err != nil {
		log.Printf("Failed to parse message: %v", err)
		return nil, err
	}

	var messageRes MessageInterface

	switch body.Type {
	case RequestMessageInfo:
		var messageInfo InfoMessage
		messageRes = &messageInfo
	case ResponseMessageInfo:
		var messageInfo InfoMessage
		messageRes = &messageInfo

	case ResponseTransactionMessage:
		var transactionMessage TransactionMessage
		messageRes = &transactionMessage
	case ResponseBlockMessage:
		var blockMessage BlockMessage
		messageRes = &blockMessage
	case ResponsePeerMessage:
		var peerMessage PeerMessage
		messageRes = &peerMessage
	case ResponseTextMessage:
		var textMessage TextMessage
		messageRes = &textMessage
	default:
		var baseMessage BaseMessage
		messageRes = &baseMessage
	}

	if err := json.Unmarshal(body.Content, &messageRes); err != nil {
		log.Printf("Failed to parse message %v: %v", body.Type, err)
		return nil, err
	}

	resultMessage := Message{
		Type:    body.Type,
		Content: messageRes,
	}
	return &resultMessage, nil

}

// Only response
func NewInfoMessage() Message {
	infoMessage := InfoMessage{
		BaseMessage: *NewBaseMessage(),
	}

	return Message{
		Type:    ResponseMessageInfo,
		Content: &infoMessage,
	}
}

func NewTransactionMessage(transaction *transaction.Transaction) Message {
	transactionMessage := TransactionMessage{
		BaseMessage: *NewBaseMessage(),
		Transaction: transaction,
	}

	return Message{
		Type:    ResponseTransactionMessage,
		Content: &transactionMessage,
	}
}

func NewBlockMessage(block *block.Block) Message {
	blockMessage := BlockMessage{
		BaseMessage: *NewBaseMessage(),
		Block:       block,
	}

	return Message{
		Type:    ResponseBlockMessage,
		Content: &blockMessage,
	}
}

func NewTextMessage(text string) Message {
	textMessage := TextMessage{
		BaseMessage: *NewBaseMessage(),
		Message:     text,
	}

	return Message{
		Type:    ResponseTextMessage,
		Content: &textMessage,
	}
}

func NewPeerMessage(ipAddr string) Message {
	peerMessage := PeerMessage{
		BaseMessage: *NewBaseMessage(),
		PeerAddrIp:  ipAddr,
	}

	return Message{
		Type:    ResponseTextMessage,
		Content: &peerMessage,
	}
}

func NewRawMessage(jsonMessage []byte) Message {
	rawMessage := RawMessage{
		BaseMessage: *NewBaseMessage(),
		MessageJson: jsonMessage,
	}

	return Message{
		Type:    RawMessageType,
		Content: &rawMessage,
	}
}
