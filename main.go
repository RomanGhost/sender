package main

import (
	"fmt"
	"log"
	"sender/data/blockchain/transaction"
	"sender/data/blockchain/wallet"
	"sender/data/deal"
	"sender/internal/server"
	"sender/internal/server/p2pprotocol"
	"sender/internal/server/p2pprotocol/serializemessage"
	"time"
)

func sendTransactions(p2p *p2pprotocol.P2PProtocol, newTransaction *transaction.Transaction) {
	for {
		time.Sleep(time.Second * 10)
		p2p.ResponseTransactionMessage(newTransaction)
		log.Println("Send transaction")

	}
}

func getDeal() *deal.Deal {
	dealJson := []byte(`{
		"id": 3,
		"buyOrder": {
			"id": 10,
			"userLogin": "roman",
			"walletID": 3,
			"cryptocurrencyCode": "BTC",
			"cardID": 2,
			"typeName": "Покупка",
			"statusName": "Используется в сделке",
			"unitPrice": 150.75,
			"quantity": 3.0,
			"description": "Созданная сделка для контрагента",
			"createdAt": "2024-11-23T01:17:14.506902",
			"lastStatusChange": "2024-11-23T01:17:14.575606"
		},
		"sellOrder": {
			"id": 2,
			"userLogin": "roman",
			"walletID": 3,
			"cryptocurrencyCode": "BTC",
			"cardID": 2,
			"typeName": "Покупка",
			"statusName": "Используется в сделке",
			"unitPrice": 150.75,
			"quantity": 3.0,
			"description": "Purchase of cryptocurrency",
			"createdAt": "2024-11-17T14:30",
			"lastStatusChange": "2024-11-23T01:24:28.737834"
		},
		"statusName": "Подтверждение сделки",
		"createdAt": "2024-11-23T01:17:14.632595",
		"lastStatusChange": "2024-11-23T01:17:14.632595"
	}`)

	dealRead, err := deal.FromJson(dealJson)
	if err != nil {
		fmt.Println(err)
		panic("Deal read error")
	}
	return dealRead
}

func main() {
	newWallet := wallet.New()
	newDeal := getDeal()

	newTransaction, _ := transaction.New(newWallet, newDeal)
	newTransaction.Sign()

	channel := make(chan serializemessage.GenericMessage)
	server := server.New("localhost", 8080, channel)
	go server.Run()
	server.Connect("localhost", 7878)

	p2pProtocol := server.GetProtocol()

	go sendTransactions(p2pProtocol, newTransaction)

	for c := range channel {

		fmt.Printf("New message: %v\n", c)
	}
}
