package serializemessage_test

import (
	"sender/internal/server/blockchain/p2pprotocol/message"
	"sender/internal/server/blockchain/p2pprotocol/message/request"
	"sender/internal/server/blockchain/p2pprotocol/message/responce"
	"sender/internal/server/blockchain/p2pprotocol/serializemessage"
	"testing"
)

func TestGenericMessage_ToJSON(t *testing.T) {
	// Создаём сообщение
	reqMessage := request.NewInfoMessage()

	// Оборачиваем в GenericMessage
	genericMessage := serializemessage.NewGenericMessage(reqMessage)

	// Сериализация
	jsonData, err := genericMessage.ToJSON()
	if err != nil {
		t.Fatalf("Failed to serialize GenericMessage: %v", err)
	}

	// Проверяем, что JSON не пустой
	if len(jsonData) == 0 {
		t.Error("Serialized JSON is empty")
	}
}

func TestGenericMessage_FromJSON(t *testing.T) {
	// Создаём оригинальное сообщение
	originalMessage := request.NewInfoMessage()

	// Оборачиваем в GenericMessage
	genericMessage := serializemessage.NewGenericMessage(originalMessage)

	// Сериализация в JSON
	jsonData, err := genericMessage.ToJSON()
	if err != nil {
		t.Fatalf("Failed to serialize GenericMessage: %v", err)
	}

	// Десериализация
	deserializedMessage, err := serializemessage.FromJSON(jsonData)
	if err != nil {
		t.Fatalf("Failed to deserialize GenericMessage: %v", err)
	}

	// Проверяем тип сообщения
	if deserializedMessage.Type != originalMessage.MessageType() {
		t.Errorf("Expected message type %s, got %s", originalMessage.MessageType(), deserializedMessage.Type)
	}

	// Проверяем контент
	deserializedContent := deserializedMessage.Content.(*request.InfoMessage)
	if deserializedContent.GetID() != originalMessage.GetID() {
		t.Errorf("Expected ID %d, got %d", originalMessage.GetID(), deserializedContent.GetID())
	}
}

func TestFromJSON_HandlesDifferentMessageTypes(t *testing.T) {
	// Тестируем десериализацию для каждого типа сообщений
	tests := []struct {
		name          string
		messageToTest message.Message
	}{
		{"BlocksBeforeMessage", request.NewBlocksBeforeMessage()},
		{"LastNBlocksMessage", request.NewLastNBlocksMessage(5)},
		{"ResponseBlockMessage", responce.NewBlockMessage(nil, false)},
		{"ResponsePeerMessage", responce.NewPeerMessage([]string{"127.0.0.1"})},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Оборачиваем в GenericMessage и сериализуем
			genericMessage := serializemessage.NewGenericMessage(tt.messageToTest)
			jsonData, err := genericMessage.ToJSON()
			if err != nil {
				t.Fatalf("Failed to serialize message %s: %v", tt.name, err)
			}

			// Десериализуем
			deserializedMessage, err := serializemessage.FromJSON(jsonData)
			if err != nil {
				t.Fatalf("Failed to deserialize message %s: %v", tt.name, err)
			}

			// Проверяем тип сообщения
			if deserializedMessage.Type != tt.messageToTest.MessageType() {
				t.Errorf("Expected message type %s, got %s", tt.messageToTest.MessageType(), deserializedMessage.Type)
			}
		})
	}
}
