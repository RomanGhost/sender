package kafkaprocessing

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaProcess struct {
	brokerAddress string
	topicName     string
	writer        *kafka.Writer
	reader        *kafka.Reader
	lastMessageId int
}

func NewKafkaProcess(brokerAddress, topicName string, numPartitions, replicationFactor int) (*KafkaProcess, error) {
	conn, err := kafka.Dial("tcp", brokerAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to broker: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return nil, fmt.Errorf("failed to get controller: %w", err)
	}

	controllerConn, err := kafka.Dial("tcp", controller.Host)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to controller: %w", err)
	}
	defer controllerConn.Close()

	topicConfig := kafka.TopicConfig{
		Topic:             topicName,
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
	}

	if err := controllerConn.CreateTopics(topicConfig); err != nil {
		return nil, fmt.Errorf("failed to create topic: %w", err)
	}

	log.Printf("Topic %s created successfully", topicName)
	return &KafkaProcess{
		brokerAddress: brokerAddress,
		topicName:     topicName,
		reader:        nil,
		writer:        nil,
		lastMessageId: 0,
	}, nil
}

func (kp *KafkaProcess) ConnectWriter() error {
	if kp.writer != nil {
		return fmt.Errorf("writer is already connected")
	}
	kp.writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kp.brokerAddress},
		Topic:    kp.topicName,
		Balancer: &kafka.LeastBytes{},
	})
	return nil
}

func (kp *KafkaProcess) CloseWriter() {
	if kp.writer != nil {
		kp.writer.Close()
		kp.writer = nil
	}
}

func (kp *KafkaProcess) ConnectReader(groupID string) error {
	if kp.reader != nil {
		return fmt.Errorf("reader is already connected")
	}
	kp.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kp.brokerAddress},
		Topic:   kp.topicName,
		GroupID: groupID,
	})
	return nil
}

func (kp *KafkaProcess) CloseReader() {
	if kp.reader != nil {
		kp.reader.Close()
		kp.reader = nil
	}
}

func (kp *KafkaProcess) WriteMessage(message string) error {
	if kp.writer == nil {
		return fmt.Errorf("Reader isn't connect")
	}

	err := kp.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(fmt.Sprintf("Key-%d", kp.lastMessageId)),
			Value: []byte(fmt.Sprintf(message)),
		},
	)
	if err != nil {
		return fmt.Errorf("write error, MessageID: %v \nError: %v", kp.lastMessageId, err)
	}
	kp.lastMessageId++
	return nil
}

func (kp *KafkaProcess) ReadMessage() (*kafka.Message, error) {
	if kp.writer == nil {
		return nil, fmt.Errorf("Writer isn't connect")
	}

	m, err := kp.reader.ReadMessage(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Error read message", err)
	}
	return &m, nil
}
