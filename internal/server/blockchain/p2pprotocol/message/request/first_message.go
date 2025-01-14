package request

import (
	"sender/internal/server/blockchain/p2pprotocol/message"
	"time"
)

type InfoMessage struct {
	message.BaseMessage
}

func NewInfoMessage() *InfoMessage {
	return &InfoMessage{
		BaseMessage: message.BaseMessage{
			ID:        0,
			Timestamp: time.Now(),
		},
	}
}

func (fm *InfoMessage) MessageType() string {
	return message.RequestMessageInfo.String()
}
