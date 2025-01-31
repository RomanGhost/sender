package process

import (
	"context"
)

type KafkaProcessInterface interface {
	WriteMessage(ctx context.Context, message string) error
	WriterConnected() bool
	GetTopicName() string
	ConnectWriter()
	CloseWriter()
}
