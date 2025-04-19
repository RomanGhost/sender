// Файл: message_test.go
package message

import (
	"encoding/json"
	"reflect"

	// Замените на реальные пути к вашим пакетам, если они существуют
	"sender/internal/data/blockchain/block"
	"sender/internal/data/blockchain/transaction"
	"sender/internal/data/blockchain/wallet"
	"sender/internal/data/deal"
	"sender/internal/data/order"
	"sync/atomic" // Используем для безопасного инкремента ID в NewBaseMessage
	"testing"
	"time"
)

// --- Тесты ---
var globalMessageID uint64

func TestMessageFromJson(t *testing.T) {
	atomic.StoreUint64(&globalMessageID, 0)
	fixedTime := time.Date(2025, time.April, 19, 19, 0, 0, 0, time.UTC).Unix()
	w := wallet.New()

	// Создаем транзакцию и связанные объекты
	buyOrder := &order.Order{ID: 1, UserHashPublicKey: w.Sereliaze().PublicKey, CryptocurrencyCode: "BTC", TypeName: "buy", UnitPrice: 50000.0, Quantity: 0.1}
	sellOrder := &order.Order{ID: 2, UserHashPublicKey: w.Sereliaze().PublicKey, CryptocurrencyCode: "BTC", TypeName: "sell", UnitPrice: 50000.0, Quantity: 0.1}
	dealObj := &deal.Deal{ID: 1, BuyOrder: buyOrder, SellOrder: sellOrder, StatusName: "completed", CreatedAt: "2025-01-01T12:00:00Z", LastStatusChange: "2025-01-01T12:30:00Z"}
	tx, _ := transaction.New(w, dealObj)
	blockObj := &block.Block{ID: 1, TimeCreated: time.Now().UTC().Unix(), Transactions: []transaction.Transaction{tx}, PreviousHash: "abc123", Nonce: 42}

	tests := []struct {
		name        string
		input       interface{}
		expectType  MessageType
		expectError bool
		expectKind  reflect.Type
	}{
		{
			name:       "Valid InfoMessage",
			input:      Message{Type: ResponseMessageInfo, Content: &InfoMessage{BaseMessage: BaseMessage{ID: 1, TimeStamp: fixedTime}}},
			expectType: ResponseMessageInfo,
			expectKind: reflect.TypeOf(&InfoMessage{}),
		},
		{
			name:       "Valid TransactionMessage",
			input:      Message{Type: ResponseTransactionMessage, Content: &TransactionMessage{BaseMessage: BaseMessage{ID: 2, TimeStamp: fixedTime}, Transaction: &tx}},
			expectType: ResponseTransactionMessage,
			expectKind: reflect.TypeOf(&TransactionMessage{}),
		},
		{
			name:       "Valid BlockMessage",
			input:      Message{Type: ResponseBlockMessage, Content: &BlockMessage{BaseMessage: BaseMessage{ID: 3, TimeStamp: fixedTime}, Block: blockObj}},
			expectType: ResponseBlockMessage,
			expectKind: reflect.TypeOf(&BlockMessage{}),
		},
		{
			name:        "Malformed JSON",
			input:       []byte(`this is not json`),
			expectError: true,
		},
		{
			name:        "Empty JSON",
			input:       []byte{},
			expectError: true,
		},
		{
			name:        "Missing content",
			input:       []byte(`{"type": "ResponseMessageInfo"}`),
			expectError: true,
		},
		{
			name:        "Invalid content field type",
			input:       []byte(`{"type": "ResponseTransactionMessage", "content": "not an object"}`),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var inputJSON []byte
			switch v := tt.input.(type) {
			case []byte:
				inputJSON = v
			default:
				inputJSON, _ = json.Marshal(v)
			}

			msg, err := MessageFromJson(inputJSON)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				if msg != nil {
					t.Errorf("Expected nil message, got: %+v", msg)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if msg.Type != tt.expectType {
				t.Errorf("Expected type %v, got %v", tt.expectType, msg.Type)
			}
			if tt.expectKind != nil && reflect.TypeOf(msg.Content) != tt.expectKind {
				t.Errorf("Expected content type %v, got %v", tt.expectKind, reflect.TypeOf(msg.Content))
			}
		})
	}
}

// --- Тесты для конструкторов New*Message ---
func TestNewInfoMessage(t *testing.T) {
	atomic.StoreUint64(&globalMessageID, 0) // Сброс ID
	startTime := time.Now().Unix()
	msg := NewInfoMessage()

	if msg.Type != ResponseMessageInfo {
		t.Errorf("Expected type %s, got %s", ResponseMessageInfo, msg.Type)
	}
	if msg.Content == nil {
		t.Fatal("Content is nil")
	}
	content, ok := msg.Content.(*InfoMessage)
	if !ok {
		t.Fatalf("Expected content type *InfoMessage, got %T", msg.Content)
	}

	if content.GetTime() < startTime {
		t.Errorf("Expected time >= %d, got %d", startTime, content.GetTime())
	}
}

func TestNewTransactionMessage(t *testing.T) {
	atomic.StoreUint64(&globalMessageID, 0) // Сброс ID
	startTime := time.Now().Unix()
	w := wallet.New()

	// Создаем заказы
	buyOrder := &order.Order{
		ID:                 1,
		UserHashPublicKey:  w.Sereliaze().PublicKey,
		CryptocurrencyCode: "BTC",
		TypeName:           "buy",
		UnitPrice:          50000.0,
		Quantity:           0.1,
	}
	sellOrder := &order.Order{
		ID:                 2,
		UserHashPublicKey:  w.Sereliaze().PublicKey,
		CryptocurrencyCode: "BTC",
		TypeName:           "sell",
		UnitPrice:          50000.0,
		Quantity:           0.1,
	}

	// Создаем сделку
	d := &deal.Deal{
		ID:               1,
		BuyOrder:         buyOrder,
		SellOrder:        sellOrder,
		StatusName:       "completed",
		CreatedAt:        "2025-01-01T12:00:00Z",
		LastStatusChange: "2025-01-01T12:30:00Z",
	}
	// Создаем транзакцию
	mockTx, _ := transaction.New(w, d)
	msg := NewTransactionMessage(&mockTx)

	if msg.Type != ResponseTransactionMessage {
		t.Errorf("Expected type %s, got %s", ResponseTransactionMessage, msg.Type)
	}
	if msg.Content == nil {
		t.Fatal("Content is nil")
	}
	content, ok := msg.Content.(*TransactionMessage)
	if !ok {
		t.Fatalf("Expected content type *TransactionMessage, got %T", msg.Content)
	}

	// Сравнение содержимого
	if !reflect.DeepEqual(content.Transaction, mockTx) {
		t.Errorf("Transaction content mismatch:\nExpected: %#v\nActual:   %#v", mockTx, content.Transaction)
	}

	if content.GetTime() < startTime {
		t.Errorf("Expected time >= %d, got %d", startTime, content.GetTime())
	}
}

func TestNewBlockMessage(t *testing.T) {
	atomic.StoreUint64(&globalMessageID, 0) // Сброс ID
	startTime := time.Now().Unix()
	mockBlock := &block.Block{} // Используем стаб/реальный тип
	msg := NewBlockMessage(mockBlock)

	if msg.Type != ResponseBlockMessage {
		t.Errorf("Expected type %s, got %s", ResponseBlockMessage, msg.Type)
	}
	if msg.Content == nil {
		t.Fatal("Content is nil")
	}
	content, ok := msg.Content.(*BlockMessage)
	if !ok {
		t.Fatalf("Expected content type *BlockMessage, got %T", msg.Content)
	}

	if content.Block != mockBlock { // Проверяем указатель
		t.Errorf("Block pointer mismatch. Expected %p, got %p", mockBlock, content.Block)
	}
	// Сравнение содержимого
	if !reflect.DeepEqual(content.Block, mockBlock) {
		t.Errorf("Block content mismatch:\nExpected: %#v\nActual:   %#v", mockBlock, content.Block)
	}

	if content.GetTime() < startTime {
		t.Errorf("Expected time >= %d, got %d", startTime, content.GetTime())
	}
}

func TestNewTextMessage(t *testing.T) {
	atomic.StoreUint64(&globalMessageID, 0) // Сброс ID
	startTime := time.Now().Unix()
	testText := "Test message content for constructor"
	msg := NewTextMessage(testText)

	if msg.Type != ResponseTextMessage {
		t.Errorf("Expected type %s, got %s", ResponseTextMessage, msg.Type)
	}
	if msg.Content == nil {
		t.Fatal("Content is nil")
	}
	content, ok := msg.Content.(*TextMessage)
	if !ok {
		t.Fatalf("Expected content type *TextMessage, got %T", msg.Content)
	}

	if content.Message != testText {
		t.Errorf("Expected text '%s', got '%s'", testText, content.Message)
	}

	if content.GetTime() < startTime {
		t.Errorf("Expected time >= %d, got %d", startTime, content.GetTime())
	}
}

func TestNewPeerMessage(t *testing.T) {
	atomic.StoreUint64(&globalMessageID, 0) // Сброс ID
	startTime := time.Now().Unix()
	testIP := "10.0.0.1:9090"
	msg := NewPeerMessage(testIP)

	if msg.Type != ResponsePeerMessage {
		t.Errorf("Expected type %s, got %s", ResponsePeerMessage, msg.Type)
	}
	if msg.Content == nil {
		t.Fatal("Content is nil")
	}
	content, ok := msg.Content.(*PeerMessage)
	if !ok {
		t.Fatalf("Expected content type *PeerMessage, got %T", msg.Content)
	}

	if content.PeerAddrIp != testIP {
		t.Errorf("Expected PeerAddrIp '%s', got '%s'", testIP, content.PeerAddrIp)
	}
	if content.GetTime() < startTime {
		t.Errorf("Expected time >= %d, got %d", startTime, content.GetTime())
	}
}

func TestNewRawMessage(t *testing.T) {
	atomic.StoreUint64(&globalMessageID, 0) // Сброс ID
	startTime := time.Now().Unix()
	testJsonData := json.RawMessage(`{"rawKey": "rawValue", "number": 123}`)
	msg := NewRawMessage(testJsonData) // Передаем json.RawMessage

	if msg.Type != RawMessageType { // Убедитесь, что RawMessageType определен
		t.Errorf("Expected type %s, got %s", RawMessageType, msg.Type)
	}
	if msg.Content == nil {
		t.Fatal("Content is nil")
	}
	content, ok := msg.Content.(*RawMessage)
	if !ok {
		t.Fatalf("Expected content type *RawMessage, got %T", msg.Content)
	}

	// Сравнение байтовых срезов (json.RawMessage - это []byte)
	if !reflect.DeepEqual(content.MessageJson, testJsonData) {
		t.Errorf("Expected MessageJson '%s', got '%s'", string(testJsonData), string(content.MessageJson))
	}

	if content.GetTime() < startTime {
		t.Errorf("Expected time >= %d, got %d", startTime, content.GetTime())
	}
}
