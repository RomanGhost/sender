package p2pprotocol

import (
	"log"
	"sender/internal/data/blockchain/block"
	"sender/internal/data/blockchain/transaction"
	"sender/internal/server/blockchain/connectionpool"
	"sender/internal/server/blockchain/p2pprotocol/message"
	"sender/internal/server/blockchain/p2pprotocol/message/request"
	"sender/internal/server/blockchain/p2pprotocol/message/responce"
	"sender/internal/server/blockchain/p2pprotocol/serializemessage"
	"sync"
)

type P2PProtocol struct {
	connectionPool connectionpool.ConnectionPoolInterface
	lastMessageID  uint64
	sender         chan<- message.Message
	mu             sync.Mutex
}

func New(connectionPool connectionpool.ConnectionPoolInterface, sender chan<- message.Message) *P2PProtocol {
	return &P2PProtocol{
		connectionPool: connectionPool,
		lastMessageID:  0,
		sender:         sender,
	}
}

func (p *P2PProtocol) HandleMessage(messageJSON string) error {
	genericMessage, err := serializemessage.FromJSON([]byte(messageJSON))
	if err != nil {
		log.Println("Failed to deserialize message:", err)
		return err
	}

	log.Printf("P2P get new message: %v\n", genericMessage.Content.MessageType())

	switch genericMessage.Type {
	case message.RequestMessageInfo.String():
		{
			p.ResponseInfoMessage()

			return nil
		}
	case message.ResponseMessageInfo.String():
		{
			messageID := genericMessage.Content.GetID()
			if p.lastMessageID < messageID {
				p.lastMessageID = messageID
			}

			return nil
		}
	}

	messageID := genericMessage.Content.GetID()
	log.Printf("P2P get new message with id: %v, My ID: %v\n", messageID, p.lastMessageID)
	if messageID <= p.lastMessageID {
		return nil
	} else {
		p.lastMessageID = messageID
	}
	log.Printf("MessageId: %v/%v\n", p.lastMessageID, messageID)

	// send message to channel
	p.sender <- genericMessage.Content
	//send everyone client
	p.Broadcast(genericMessage.Content, true)
	return nil
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

func (p *P2PProtocol) ResponcePeerMessage() {
	addresses := p.connectionPool.GetPeerAddresses()
	message := responce.NewPeerMessage(addresses)

	p.Broadcast(message, false)
}

func (p *P2PProtocol) GetLastMessageId() uint64 {
	return p.lastMessageID
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
