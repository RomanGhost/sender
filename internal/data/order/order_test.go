package order_test

import (
	"sender/internal/data/order"
	"strings"
	"testing"
)

func TestOrderToJson(t *testing.T) {
	// Подготовка тестовых данных
	testOrder := &order.Order{
		ID:                 1,
		UserHashPublicKey:  "user_public_key",
		CryptocurrencyCode: "BTC",
		TypeName:           "buy",
		UnitPrice:          50000.50,
		Quantity:           2.5,
		CreatedAt:          "2023-01-01T10:00:00Z",
		LastStatusChange:   "2023-01-01T12:00:00Z",
	}

	// Тестирование метода ToJson
	jsonData, err := testOrder.ToJson()
	if err != nil {
		t.Fatalf("Failed to convert order to JSON: %v", err)
	}

	// Проверка: JSON должен содержать ключи и значения из заказа
	jsonString := string(jsonData)
	expectedStrings := []string{
		`"id":1`,
		`"userHashPublicKey":"user_public_key"`,
		`"cryptocurrencyCode":"BTC"`,
		`"typeName":"buy"`,
		`"unitPrice":50000.5`,
		`"quantity":2.5`,
		`"createdAt":"2023-01-01T10:00:00Z"`,
		`"lastStatusChange":"2023-01-01T12:00:00Z"`,
	}
	for _, expected := range expectedStrings {
		if !strings.Contains(jsonString, expected) {
			t.Errorf("JSON representation is incorrect, missing: %s", expected)
		}
	}
}

func TestOrderFromJson(t *testing.T) {
	// Подготовка JSON-данных
	jsonData := `{
		"id": 1,
		"userHashPublicKey": "user_public_key",
		"cryptocurrencyCode": "BTC",
		"typeName": "buy",
		"unitPrice": 50000.5,
		"quantity": 2.5,
		"createdAt": "2023-01-01T10:00:00Z",
		"lastStatusChange": "2023-01-01T12:00:00Z"
	}`

	// Тестирование метода FromJson
	order, err := order.FromJson([]byte(jsonData))
	if err != nil {
		t.Fatalf("Failed to parse JSON to order: %v", err)
	}

	// Проверка полей объекта Order
	if order.ID != 1 {
		t.Errorf("Expected ID 1, got %d", order.ID)
	}
	if order.UserHashPublicKey != "user_public_key" {
		t.Errorf("Expected UserHashPublicKey 'user_public_key', got %s", order.UserHashPublicKey)
	}
	if order.CryptocurrencyCode != "BTC" {
		t.Errorf("Expected CryptocurrencyCode 'BTC', got %s", order.CryptocurrencyCode)
	}
	if order.TypeName != "buy" {
		t.Errorf("Expected TypeName 'buy', got %s", order.TypeName)
	}
	if order.UnitPrice != 50000.5 {
		t.Errorf("Expected UnitPrice 50000.5, got %f", order.UnitPrice)
	}
	if order.Quantity != 2.5 {
		t.Errorf("Expected Quantity 2.5, got %f", order.Quantity)
	}
	if order.CreatedAt != "2023-01-01T10:00:00Z" {
		t.Errorf("Expected CreatedAt '2023-01-01T10:00:00Z', got %s", order.CreatedAt)
	}
	if order.LastStatusChange != "2023-01-01T12:00:00Z" {
		t.Errorf("Expected LastStatusChange '2023-01-01T12:00:00Z', got %s", order.LastStatusChange)
	}
}
