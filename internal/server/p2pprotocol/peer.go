package p2pprotocol

import (
	"fmt"
	"log"
	"sender/data/blockchain/block"
	"sender/data/blockchain/transaction"
	"sender/internal/server/connectionpool"
	"sender/internal/server/p2pprotocol/message"
	"sender/internal/server/p2pprotocol/message/request"
	"sender/internal/server/p2pprotocol/message/responce"
	"sender/internal/server/p2pprotocol/serializemessage"
	"sync"
)

type P2PProtocol struct {
	connectionPool *connectionpool.ConnectionPool
	lastMessageID  uint64
	sender         chan<- serializemessage.GenericMessage
	mu             sync.Mutex
}

func New(connectionPool *connectionpool.ConnectionPool, sender chan<- serializemessage.GenericMessage) *P2PProtocol {
	return &P2PProtocol{
		connectionPool: connectionPool,
		lastMessageID:  0,
		sender:         sender,
	}
}

func (p *P2PProtocol) HandleMessage(messageJSON string) {
	genericMessage, err := serializemessage.FromJSON([]byte(messageJSON))
	if err != nil {
		log.Println("Failed to deserialize message:", err)
		return
	}

	log.Printf("P2P get new message: %v\n", genericMessage.Content.MessageType())

	switch genericMessage.Type {
	case "RequestMessageInfo":
		{
			p.ResponseInfoMessage()

			return
		}
	case "ResponseMessageInfo":
		{
			messageID := genericMessage.Content.GetID()
			if p.lastMessageID < messageID {
				p.lastMessageID = messageID
			}

			return
		}
	}

	messageID := genericMessage.Content.GetID()
	log.Printf("P2P get new message with id: %v, My ID: %v\n", messageID, p.lastMessageID)
	if messageID <= p.lastMessageID {
		return
	} else {
		p.lastMessageID = messageID
	}
	fmt.Printf("MessageId: %v/%v\n", p.lastMessageID, messageID)

	// send message to channel
	p.sender <- *genericMessage
	//send everyone client
	p.Broadcast(genericMessage.Content, true)
}

func (p *P2PProtocol) RequestInfoMessage() {
	message := request.NewInfoMessage()

	p.Broadcast(message, false)
}

func (p *P2PProtocol) RequestChainMessage(countBlocks uint) {
	message := request.NewLastNBlocksMessage(countBlocks)

	p.Broadcast(message, false)
}

func (p *P2PProtocol) RequestBlocksBeforeMessage() {
	message := request.NewBlocksBeforeMessage()

	p.Broadcast(message, false)
}

func (p *P2PProtocol) ResponseInfoMessage() {
	message := responce.NewInfoMessage()

	p.Broadcast(message, false)
}

func (p *P2PProtocol) ResponseBlockMessage(sendBlock *block.Block) {
	message := responce.NewBlockMessage(sendBlock, false)

	p.Broadcast(message, false)
}

func (p *P2PProtocol) ResponseTransactionMessage(sendTransaction *transaction.Transaction) {
	message := responce.NewTransactionMessage(sendTransaction)

	p.Broadcast(message, false)
}

func (p *P2PProtocol) ResponseChainMessage(sendChain []block.Block) {
	message := responce.NewChainMessage(sendChain)

	p.Broadcast(message, false)
}

func (p *P2PProtocol) Broadcast(getMessage message.Message, receive bool) {
	p.mu.Lock()
	if !receive {
		p.lastMessageID++
	}
	getMessage.SetID(p.lastMessageID)
	p.mu.Unlock()

	genericMessage := serializemessage.NewGenericMessage(getMessage)
	jsonText, _ := genericMessage.ToJSON()
	jsonString := string(jsonText)

	p.connectionPool.Broadcast(jsonString)
}
