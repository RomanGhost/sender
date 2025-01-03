package main

import (
	"context"
	"fmt"
	"log"
	"sender/internal/data/blockchain/transaction"
	"sender/internal/data/blockchain/wallet"
	"sender/internal/data/deal"
	"sender/internal/process"
	"sender/internal/server"
	"sender/internal/server/p2pprotocol"
	"sender/internal/server/p2pprotocol/message"
	"sync"
	"time"
)

func sendTransactions(p2p *p2pprotocol.P2PProtocol, newTransaction *transaction.Transaction) {
	for {
		time.Sleep(time.Second * 30)
		p2p.ResponseTransactionMessage(newTransaction)
		log.Println("Send transaction")

	}
}

func getDeal() *deal.Deal {
	dealJson := []byte(`{
		"id": 3,
		"buyOrder": {
			"id": 10,
			"userHashPublicKey": "sha256",
			"cryptocurrencyCode": "BTC",
			"typeName": "Покупка",
			"unitPrice": 150.75,
			"quantity": 3.0,
			"createdAt": "2024-11-23T01:17:14.506902",
			"lastStatusChange": "2024-11-23T01:17:14.575606"
		},
		"sellOrder": {
			"id": 2,
			"userHashPublicKey": "sha256",
			"cryptocurrencyCode": "BTC",
			"typeName": "Продажа",
			"unitPrice": 150.75,
			"quantity": 3.0,
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

func writeKafkaMessage(kafka_process *process.KafkaProcess, start, end int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := start; i < end; i++ {
		fmt.Println(i)
		kafka_process.WriteMessage(context.Background(), fmt.Sprintf("Hello! %v", i))
	}
}

func readKafkaMessage(kafka_process *process.KafkaProcess) {
	kafka_process.ConnectReader()
	defer kafka_process.CloseReader()

	handleMessage := func(msg string) {
		log.Printf("Processing message: %s", msg)
	}

	err := kafka_process.ReadMessages(context.Background(), handleMessage)
	if err != nil {
		log.Fatal("Error of reading", err)
	}
}

func main() {
	kafka_process_consumer := process.NewKafkaProcess("localhost:9092", "GoGetDeal", "middle-group")
	go readKafkaMessage(kafka_process_consumer)

	newWallet := wallet.New()
	newDeal := getDeal()

	newTransaction, _ := transaction.New(newWallet, newDeal)
	newTransaction.Sign()

	channel := make(chan message.Message)
	defer close(channel)

	serverBlockchain := server.New("localhost", 7990, channel)
	go serverBlockchain.Run()
	err := serverBlockchain.Connect("localhost", 7878)
	if err != nil {
		fmt.Printf("Coudn't connect to server: %v\n", err)
	}

	p2pProtocol := serverBlockchain.GetProtocol()

	go sendTransactions(p2pProtocol, newTransaction)

	if err != nil {
		log.Fatalln("Error with topic kafka: ", err)
	}

	kafka_process_producer := process.NewKafkaProcess("localhost:9092", "SpringGetDeal", "example-group")
	kafka_process_producer.ConnectWriter()
	defer kafka_process_producer.Close()

	go process.MessageProcessing(channel, p2pProtocol, kafka_process_producer)

	fmt.Println("Enter q to exit: ")
	for {
		text2 := ""
		fmt.Scanln(text2)
		if text2 == "q" {
			break
		}
	}
}
