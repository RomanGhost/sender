package app

import (
	"fmt"
	"log"
	"sender/internal/data/blockchain/transaction"
	"sender/internal/server/blockchain"
	"sender/internal/server/blockchain/protocol/message"
)

type AppState struct {
	Server       *blockchain.Server
	KafkaChan    chan message.MessageInterface
	ProtocolChan chan message.Message
}

// func NewAppState(server *blockchain.Server) AppState {
// 	return AppState{
// 		Server:          server,
// 		KafkaChannel:    make(chan message.MessageInterface),
// 		ProtocolChannel: make(chan message.Message),
// 	}
// }

func (s *AppState) ReadBlock(blockMessage *message.BlockMessage) {
	s.KafkaChan <- blockMessage
}

func (s *AppState) SendTransaction(transaction *transaction.Transaction) {
	messageTransaction := message.NewTransactionMessage(transaction)
	s.ProtocolChan <- messageTransaction
}

func (s *AppState) Connect(ipAddr string) {
	err := s.Server.Connect(fmt.Sprintf("%v:7878", ipAddr))
	if err != nil {
		log.Printf("Couldn't connect by addr: %v, err: %v", ipAddr, err)
	}
}
