package request

import (
	"sender/internal/server/p2pprotocol/message"
	"time"
)

type FirstMessage struct {
	message.BaseMessage
}

func NewFirstMessage() *FirstMessage {
	return &FirstMessage{
		BaseMessage: message.BaseMessage{
			ID:        0,
			Timestamp: time.Now(),
		},
	}
}

func (fm *FirstMessage) MessageType() string {
	return "ResponseMessageInfo"
}
