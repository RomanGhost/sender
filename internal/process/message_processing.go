package process

import (
	"context"
	"fmt"
	"log"
	"sender/internal/server/p2pprotocol"
	"sender/internal/server/p2pprotocol/message"
	"sender/internal/server/p2pprotocol/message/responce"
)

func MessageProcessing(c chan message.Message, p2pProtocol *p2pprotocol.P2PProtocol, kafkaWriter *KafkaProcess) error {
	if !kafkaWriter.WriterConnected() {
		return fmt.Errorf("Writer isn't open")
	}

	for messageGet := range c {
		switch msg := messageGet.(type) {
		case *responce.BlockMessage:
			for _, tr := range msg.Block.Transactions {
				log.Printf("Transaction send to kafka topic: %s", kafkaWriter.TopicName)
				message := tr.DealMessage
				kafkaWriter.WriteMessage(context.Background(), string(message))
			}
		case *responce.PeerMessage:
			for _, pa := range msg.PeerAddresses {
				fmt.Printf("Connection addresses: %v\n", pa)
			}
		default:
			fmt.Println("Неизвестный тип сообщения")
		}
	}
	return nil
}
