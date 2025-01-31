package process

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type WriterInterface interface {
	WriteMessages(ctx context.Context, messages ...kafka.Message) error
	Close() error
}

type ReaderInterface interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	Close() error
}

type KafkaProcess struct {
	BrokerAddress string
	TopicName     string
	Reader        ReaderInterface
	Writer        WriterInterface
	GroupID       string
}

// NewKafkaProcess initializes a KafkaProcess with provided configuration.
func NewKafkaProcess(brokerAddress, topicName, groupID string) *KafkaProcess {
	return &KafkaProcess{
		BrokerAddress: brokerAddress,
		TopicName:     topicName,
		GroupID:       groupID,
	}
}

// ConnectWriter sets up a Kafka writer.
func (kp *KafkaProcess) ConnectWriter() {
	kp.Writer = &kafka.Writer{
		Addr:         kafka.TCP(kp.BrokerAddress),
		Topic:        kp.TopicName,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
	}
	log.Println("Kafka writer connected")
}

func (kp *KafkaProcess) GetTopicName() string {
	return kp.TopicName
}

func (kp *KafkaProcess) WriterConnected() bool {
	return kp.Writer != nil
}

func (kp *KafkaProcess) CloseWriter() {
	if kp.Writer != nil {
		kp.Writer.Close()
		kp.Writer = nil
	}
}

func (kp *KafkaProcess) ConnectReader() {
	kp.Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kp.BrokerAddress},
		Topic:    kp.TopicName,
		GroupID:  kp.GroupID,
		MinBytes: 10e3, // 10 KB
		MaxBytes: 10e6, // 10 MB
		MaxWait:  1 * time.Second,
	})
	log.Println("Kafka reader connected")
}

func (kp *KafkaProcess) CloseReader() {
	if kp.Reader != nil {
		kp.Reader.Close()
		kp.Reader = nil
	}
}

func (kp *KafkaProcess) ReaderConnected() bool {
	return kp.Reader != nil
}

func (kp *KafkaProcess) Close() {
	if kp.Writer != nil {
		kp.CloseWriter()
	}
	if kp.Reader != nil {
		kp.CloseReader()
	}
	log.Println("Kafka connections closed")
}

// WriteMessage sends a message to the Kafka topic.
func (kp *KafkaProcess) WriteMessage(ctx context.Context, message string) error {
	if kp.Writer == nil {
		return errors.New("Kafka writer is not initialized")
	}

	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("DeaiAdd-%d", time.Now().UnixNano())),
		Value: []byte(message),
	}

	err := kp.Writer.WriteMessages(ctx, msg)
	if err != nil {
		log.Printf("Failed to write message: %v", err)
		return err
	}

	log.Printf("Message written: %s", message)
	return nil
}

// ReadMessages listens for messages from the Kafka topic.
func (kp *KafkaProcess) ReadMessages(ctx context.Context, handleMessage func(string)) error {
	if kp.Reader == nil {
		return errors.New("Kafka reader is not initialized")
	}

	for {
		msg, err := kp.Reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				log.Println("Reader context canceled")
				return nil
			}
			log.Printf("Failed to read message: %v", err)
			continue
		}

		message := string(msg.Value)
		log.Printf("Message received from kafka: %s", message)

		// Process the message
		handleMessage(message)
	}
}
