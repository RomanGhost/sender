package p2pprotocol

import (
	"encoding/json"
	"log"
	"sender/internal/p2pprotocol/message"
	"sender/internal/server/connectionpool"
	"sync"
)

type P2PProtocol struct {
	connectionPool *connectionpool.ConnectionPool
	lastMessageID  uint64
	sender         chan message.Message
	mu             sync.Mutex
}

func New(connectionPool *connectionpool.ConnectionPool, sender chan message.Message) *P2PProtocol {
	return &P2PProtocol{
		connectionPool: connectionPool,
		lastMessageID:  0,
		sender:         sender,
	}
}

func (p *P2PProtocol) HandleMessage(messageJSON string) {
	var genericMessage message.GenericMessage
	messageBase := genericMessage.Content
	err := json.Unmarshal([]byte(messageJSON), &messageBase)
	if err != nil {
		log.Println("Failed to deserialize message:", err)
		return
	}

	if messageBase.GetID() <= p.lastMessageID {
		return
	}

	p.lastMessageID = messageBase.GetID()

	// Пример обработки конкретного типа сообщения.
	if messageBase.GetID() == 0 {
		log.Println("Processing RequestMessageInfo")
	}
	// Отправка сообщения в канал.
	p.sender <- messageBase
	p.Broadcast(messageBase, true)
}

// func (p *P2PProtocol) ResponseFirstMessage() {
// 	message := &RequestMessageInfo{}
// 	message.SetID(p.lastMessageID + 1)
// 	p.Broadcast(message, false)
// }

func (p *P2PProtocol) Broadcast(getMessage message.Message, receive bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !receive {
		p.lastMessageID++
	}
	getMessage.SetID(p.lastMessageID)

	genericMessage := message.NewGenericMessage(getMessage)
	jsonText, _ := genericMessage.ToJSON()
	jsonString := string(jsonText)

	p.connectionPool.Broadcast(jsonString)
}
