package process

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWriter implements WriterInterface for testing
type MockWriter struct {
	mock.Mock
}

func (m *MockWriter) WriteMessages(ctx context.Context, messages ...kafka.Message) error {
	args := m.Called(ctx, messages)
	return args.Error(0)
}

func (m *MockWriter) Close() error {
	args := m.Called()
	return args.Error(0)
}

// MockReader implements ReaderInterface for testing
type MockReader struct {
	mock.Mock
}

func (m *MockReader) ReadMessage(ctx context.Context) (kafka.Message, error) {
	args := m.Called(ctx)
	return args.Get(0).(kafka.Message), args.Error(1)
}

func (m *MockReader) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestNewKafkaProcess(t *testing.T) {
	// Arrange
	brokerAddress := "localhost:9092"
	topicName := "test-topic"
	groupID := "test-group"

	// Act
	kp := NewKafkaProcess(brokerAddress, topicName, groupID)

	// Assert
	assert.NotNil(t, kp)
	assert.Equal(t, brokerAddress, kp.BrokerAddress)
	assert.Equal(t, topicName, kp.TopicName)
	assert.Equal(t, groupID, kp.GroupID)
	assert.Nil(t, kp.Reader)
	assert.Nil(t, kp.Writer)
}

func TestWriterConnected(t *testing.T) {
	// Arrange
	kp := NewKafkaProcess("localhost:9092", "test-topic", "test-group")
	mockWriter := new(MockWriter)

	// Test when writer is not connected
	assert.False(t, kp.WriterConnected())

	// Test when writer is connected
	kp.Writer = mockWriter
	assert.True(t, kp.WriterConnected())
}

func TestReaderConnected(t *testing.T) {
	// Arrange
	kp := NewKafkaProcess("localhost:9092", "test-topic", "test-group")
	mockReader := new(MockReader)

	// Test when reader is not connected
	assert.False(t, kp.ReaderConnected())

	// Test when reader is connected
	kp.Reader = mockReader
	assert.True(t, kp.ReaderConnected())
}

func TestWriteMessage(t *testing.T) {
	// Arrange
	kp := NewKafkaProcess("localhost:9092", "test-topic", "test-group")
	mockWriter := new(MockWriter)
	ctx := context.Background()
	message := "test message"

	t.Run("Writer not initialized", func(t *testing.T) {
		// Act
		err := kp.WriteMessage(ctx, message)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, "Kafka writer is not initialized", err.Error())
	})

	t.Run("Successful write", func(t *testing.T) {
		// Arrange
		kp.Writer = mockWriter
		mockWriter.On("WriteMessages", ctx, mock.Anything).Return(nil)

		// Act
		err := kp.WriteMessage(ctx, message)

		// Assert
		assert.NoError(t, err)
		mockWriter.AssertExpectations(t)
	})

	t.Run("Write error", func(t *testing.T) {
		// Arrange
		expectedErr := errors.New("Kafka writer is not initialized")
		mockWriter.On("WriteMessages", ctx, mock.Anything).Return(expectedErr)

		// Act
		kp.Writer = nil
		err := kp.WriteMessage(ctx, message)
		fmt.Println("Error: ", err)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestReadMessages(t *testing.T) {
	// Arrange
	kp := NewKafkaProcess("localhost:9092", "test-topic", "test-group")
	mockReader := new(MockReader)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	t.Run("Reader not initialized", func(t *testing.T) {
		// Act
		err := kp.ReadMessages(ctx, func(s string) {})

		// Assert
		assert.Error(t, err)
		assert.Equal(t, "Kafka reader is not initialized", err.Error())
	})

	t.Run("Successful read", func(t *testing.T) {
		// Arrange
		kp.Reader = mockReader
		expectedMsg := "test message"
		messageReceived := false

		mockReader.On("ReadMessage", ctx).Return(
			kafka.Message{Value: []byte(expectedMsg)},
			nil,
		).Once()

		mockReader.On("ReadMessage", ctx).Return(
			kafka.Message{},
			context.Canceled,
		)

		// Act
		err := kp.ReadMessages(ctx, func(s string) {
			messageReceived = true
			assert.Equal(t, expectedMsg, s)
		})

		// Assert
		assert.NoError(t, err)
		assert.True(t, messageReceived)
		mockReader.AssertExpectations(t)
	})

	t.Run("Read error and continue", func(t *testing.T) {
		// Arrange
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		kp.Reader = mockReader
		defer cancel()

		mockReader.On("ReadMessage", ctx).Return(
			kafka.Message{},
			errors.New("read error"),
		).Once()

		mockReader.On("ReadMessage", ctx).Return(
			kafka.Message{},
			context.Canceled,
		)

		// Act
		err := kp.ReadMessages(ctx, func(s string) {})

		// Assert
		assert.NoError(t, err)
		mockReader.AssertExpectations(t)
	})
}

func TestClose(t *testing.T) {
	// Arrange
	kp := NewKafkaProcess("localhost:9092", "test-topic", "test-group")
	mockReader := new(MockReader)
	mockWriter := new(MockWriter)

	mockReader.On("Close").Return(nil)
	mockWriter.On("Close").Return(nil)

	kp.Reader = mockReader
	kp.Writer = mockWriter

	// Act
	kp.Close()

	// Assert
	assert.Nil(t, kp.Reader)
	assert.Nil(t, kp.Writer)
	mockReader.AssertExpectations(t)
	mockWriter.AssertExpectations(t)
}
