package main

import (
	"fmt"
	"log"
	"sender/internal/data/blockchain/transaction"
	"sender/internal/data/deal"
	kafkaprocessing "sender/internal/process/kafka_processing"
	"sender/internal/server/p2pprotocol"
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

func writeKafkaMessage(kafka_process *kafkaprocessing.KafkaProcess, start, end int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := start; i < end; i++ {
		fmt.Println(i)
		kafka_process.WriteMessage(fmt.Sprintf("Hello! %v", i))
	}
}

func readKafkaMessage(kafka_process *kafkaprocessing.KafkaProcess) {
	kafka_process.ConnectReader("example-group")
	defer kafka_process.CloseWriter()

	for {
		message, err := kafka_process.ReadMessage()
		if err != nil {
			log.Fatal("Error of reading", err)
		}
		log.Printf("Read new message from kafka: %s", message.Value)
	}
}

func main() {
	// kafka_process, err := kafkaprocessing.NewKafkaProcess("localhost:9092", "GetDeal", 5, 3)
	// if err != nil {
	// 	log.Fatalln("Error with topic kafka: ", err)
	// }
	// kafka_process.ConnectWriter()
	// defer kafka_process.CloseWriter()

	// go readKafkaMessage(kafka_process)

	// var wg sync.WaitGroup
	// n := 10
	// start := 0
	// for i := range 10 {
	// 	fmt.Println(i)
	// 	go writeKafkaMessage(kafka_process, start, start+n, &wg)
	// 	start += n
	// 	wg.Add(1)
	// }
	// wg.Wait()

	// newWallet := wallet.New()
	// newDeal := getDeal()

	// newTransaction, _ := transaction.New(newWallet, newDeal)
	// newTransaction.Sign()

	// channel := make(chan message.Message)
	// defer close(channel)

	// serverBlockchain := server.New("localhost", 8080, channel)
	// go serverBlockchain.Run()
	// err = serverBlockchain.Connect("localhost", 7878)
	// if err != nil {
	// 	fmt.Printf("Coudn't connect to server: %v\n", err)
	// }

	// p2pProtocol := serverBlockchain.GetProtocol()

	// go sendTransactions(p2pProtocol, newTransaction)
	// go process.MessageProcessing(channel, p2pProtocol)

	// fmt.Println("Enter q to exit: ")
	// for {
	// 	text2 := ""
	// 	fmt.Scanln(text2)
	// 	if text2 == "q" {
	// 		break
	// 	}
	// }
}
