package process

import (
	"fmt"
	"sender/internal/server/p2pprotocol"
	"sender/internal/server/p2pprotocol/message"
	"sender/internal/server/p2pprotocol/message/responce"
)

func MessageProcessing(c chan message.Message, p2pProtocol *p2pprotocol.P2PProtocol) {

	for messageGet := range c {

		switch msg := messageGet.(type) { // Type assertion with type switch
		case *responce.BlockMessage:
			for _, tr := range msg.Block.Transactions {
				fmt.Printf("Transaction from block: %v\n", tr.Transfer)
			}
		case *responce.PeerMessage:
			for _, pa := range msg.PeerAddresses {
				fmt.Printf("Connection addresses: %v\n", pa)
			}
		default:
			fmt.Println("Неизвестный тип сообщения")
		}
	}
}
