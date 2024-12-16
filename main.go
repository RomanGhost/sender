package main

import (
	"fmt"
	"sender/internal/data/blockchain/transaction"
	"sender/internal/data/blockchain/wallet"
	"sender/internal/data/deal"
	"sender/internal/p2pprotocol/message"
	"sender/internal/p2pprotocol/message/responce"
	"sender/internal/server"
	"time"
)

func getDeal() *deal.Deal {
	dealJson := []byte(`{
		"id": 3,
		"buyOrder": {
			"id": 10,
			"userLogin": "roman",
			"walletId": 3,
			"cryptocurrencyCode": "BTC",
			"cardId": 2,
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
			"walletId": 3,
			"cryptocurrencyCode": "BTC",
			"cardId": 2,
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
	transactionMessage := responce.NewTransactionMessage(newTransaction)

	channel := make(chan message.Message)
	server := server.New("localhost", 8080, channel)
	go server.Run()
	server.Connect("localhost", 7878)

	p2pProtocol := server.GetProtocol()

	time.Sleep(5 * time.Second) // Останавливает выполнение main на 30 секунд
	p2pProtocol.Broadcast(transactionMessage, false)

	time.Sleep(290 * time.Second)
	fmt.Println("Код успешно завершается!")
	for {
	}
}
