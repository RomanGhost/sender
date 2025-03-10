package main

import (
	"context"
	"fmt"
	"log"
	"sender/internal/data/blockchain/transaction"
	"sender/internal/data/blockchain/wallet"
	"sender/internal/data/deal"
	"sender/internal/process"
	"sender/internal/server/blockchain"
	"sender/internal/server/blockchain/p2pprotocol"
	"sender/internal/server/blockchain/p2pprotocol/message"
	"sender/internal/server/web"
	"sync"
)

func readKafkaMessage(kafka_process *process.KafkaProcess, p2pProtocol *p2pprotocol.P2PProtocol, wallet *wallet.Wallet) {
	kafka_process.ConnectReader()
	defer kafka_process.CloseReader()

	handleMessage := func(msg string) {
		log.Printf("Processing message from kafka: %s", msg)

		newDeal, err := deal.FromJson([]byte(msg))
		if err != nil {
			fmt.Println(err)
			panic("Deal read error")
		}

		newTransaction, _ := transaction.New(wallet, newDeal)
		newTransaction.Sign()

		p2pProtocol.ResponseTransactionMessage(newTransaction)
	}

	err := kafka_process.ReadMessages(context.Background(), handleMessage)
	if err != nil {
		log.Fatal("Error of reading", err)
	}
}

func main() {
	var wg sync.WaitGroup
	channel := make(chan message.Message)
	defer close(channel)

	//Blockchain
	serverBlockchain := blockchain.New("0.0.0.0", 7990, channel)
	wg.Add(1)
	go serverBlockchain.Run()

	err := serverBlockchain.Connect("172.17.0.2", 7878)
	if err != nil {
		fmt.Printf("Coudn't connect to server: %v\n", err)
	}

	p2pProtocol := serverBlockchain.GetProtocol()

	// go sendTransactions(p2pProtocol, newTransaction)

	// if err != nil {
	// 	log.Fatalln("Error with topic kafka: ", err)
	// }

	// Kafka connect

	kafka_process_producer := process.NewKafkaProcess("localhost:9092", "SpringGetDeal", "example-group")
	kafka_process_producer.ConnectWriter()
	defer kafka_process_producer.Close()

	wg.Add(1)
	go process.MessageProcessing(channel, p2pProtocol, kafka_process_producer)

	kafka_process_consumer := process.NewKafkaProcess("localhost:9092", "GoGetDeal", "middle-group")
	newWallet := wallet.New()

	wg.Add(1)
	go readKafkaMessage(kafka_process_consumer, p2pProtocol, newWallet)

	// web server setting
	web_server := web.New("7980")
	wg.Add(1)
	go web_server.Run()

	wg.Wait()
}
