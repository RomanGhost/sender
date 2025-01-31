package process

import (
	"context"
	"errors"
	"sender/internal/data/blockchain/block"
	"sender/internal/server/blockchain/p2pprotocol"
	"sender/internal/server/blockchain/p2pprotocol/message"
	"sender/internal/server/blockchain/p2pprotocol/message/responce"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockKafkaProcess mocks KafkaProcess for testing
type MockKafkaProcess struct {
	mock.Mock
	TopicName string
	writer    WriterInterface
}

// GetTopicName implements KafkaProcessInterface.
func (m *MockKafkaProcess) GetTopicName() string {
	panic("unimplemented")
}

func (m *MockKafkaProcess) WriteMessage(ctx context.Context, message string) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *MockKafkaProcess) WriterConnected() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockKafkaProcess) ConnectWriter() {
	return
}

func (m *MockKafkaProcess) CloseWriter() {
	return
}

// MockBlock represents a mock block for testing
type MockBlock struct {
	Transactions []MockTransaction
}

type MockTransaction struct {
	DealMessage []byte
}

func TestMessageProcessing(t *testing.T) {
	t.Run("Writer not connected", func(t *testing.T) {
		// Arrange
		mockKafka := new(MockKafkaProcess)
		mockKafka.On("WriterConnected").Return(false)

		c := make(chan message.Message)
		p2p := &p2pprotocol.P2PProtocol{}

		// Act
		err := MessageProcessing(c, p2p, mockKafka)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, "Writer isn't open", err.Error())

	})

	t.Run("Process BlockMessage successfully", func(t *testing.T) {
		// Arrange
		mockKafka := new(MockKafkaProcess)
		mockKafka.On("WriterConnected").Return(true)
		mockKafka.TopicName = "test-topic"

		mockKafka.On("WriteMessage", mock.Anything, "deal1").Return(nil)
		mockKafka.On("WriteMessage", mock.Anything, "deal2").Return(nil)

		c := make(chan message.Message)
		p2p := &p2pprotocol.P2PProtocol{}

		// Start processing in goroutine
		go func() {
			blockMsg := responce.NewBlockMessage(&block.Block{}, false)
			c <- blockMsg
			close(c)
		}()

		// Act
		err := MessageProcessing(c, p2p, mockKafka)

		// Assert
		assert.NoError(t, err)
		//
	})

	t.Run("Process PeerMessage successfully", func(t *testing.T) {
		// Arrange
		mockKafka := new(MockKafkaProcess)
		mockKafka.On("WriterConnected").Return(true)

		c := make(chan message.Message)
		p2p := &p2pprotocol.P2PProtocol{}

		// Start processing in goroutine
		go func() {
			peerMsg := &responce.PeerMessage{
				PeerAddresses: []string{"peer1", "peer2"},
			}
			c <- peerMsg
			close(c)
		}()

		// Act
		err := MessageProcessing(c, p2p, mockKafka)

		// Assert
		assert.NoError(t, err)

	})

	t.Run("Process unknown message type", func(t *testing.T) {
		// Arrange
		mockKafka := new(MockKafkaProcess)
		mockKafka.On("WriterConnected").Return(true)

		c := make(chan message.Message)
		p2p := &p2pprotocol.P2PProtocol{}

		// Start processing in goroutine
		go func() {
			// Send an unknown message type
			c <- struct{ message.Message }{}
			close(c)
		}()

		// Act
		err := MessageProcessing(c, p2p, mockKafka)

		// Assert
		assert.NoError(t, err)

	})

	t.Run("Handle kafka write error", func(t *testing.T) {
		// Arrange
		mockKafka := new(MockKafkaProcess)
		mockKafka.On("WriterConnected").Return(true)
		mockKafka.TopicName = "test-topic"

		expectedErr := errors.New("kafka write error")
		mockKafka.On("WriteMessage", mock.Anything, "deal1").Return(expectedErr)

		c := make(chan message.Message)
		p2p := &p2pprotocol.P2PProtocol{}

		// Start processing in goroutine
		go func() {
			blockMsg := responce.NewBlockMessage(&block.Block{}, false)
			c <- blockMsg
			close(c)
		}()

		// Act
		err := MessageProcessing(c, p2p, mockKafka)

		// Assert
		assert.NoError(t, err) // The function continues processing despite write errors
		//
	})
}

func TestMessageProcessingChannelClosed(t *testing.T) {
	// Arrange
	mockKafka := new(MockKafkaProcess)
	mockKafka.On("WriterConnected").Return(true)

	c := make(chan message.Message)
	p2p := &p2pprotocol.P2PProtocol{}

	// Close channel immediately
	close(c)

	// Act
	err := MessageProcessing(c, p2p, mockKafka)

	// Assert
	assert.NoError(t, err)

}
