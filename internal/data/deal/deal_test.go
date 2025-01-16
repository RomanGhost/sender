package deal_test

import (
	"sender/internal/data/deal"
	"sender/internal/data/order"
	"strings"
	"testing"
)

func TestDealToJson(t *testing.T) {
	// Подготовка тестовых данных
	buyOrder := &order.Order{
		ID:                 1,
		UserHashPublicKey:  "buyer_public_key",
		CryptocurrencyCode: "BTC",
		TypeName:           "buy",
		UnitPrice:          50000.5,
		Quantity:           1.0,
		CreatedAt:          "2023-01-01T10:00:00Z",
		LastStatusChange:   "2023-01-01T10:05:00Z",
	}

	sellOrder := &order.Order{
		ID:                 2,
		UserHashPublicKey:  "seller_public_key",
		CryptocurrencyCode: "BTC",
		TypeName:           "sell",
		UnitPrice:          50000.5,
		Quantity:           1.0,
		CreatedAt:          "2023-01-01T09:50:00Z",
		LastStatusChange:   "2023-01-01T10:00:00Z",
	}

	testDeal := &deal.Deal{
		ID:               1,
		BuyOrder:         buyOrder,
		SellOrder:        sellOrder,
		StatusName:       "completed",
		CreatedAt:        "2023-01-01T10:10:00Z",
		LastStatusChange: "2023-01-01T10:15:00Z",
	}

	// Тестирование метода ToJson
	jsonData, err := testDeal.ToJson()
	if err != nil {
		t.Fatalf("Failed to convert deal to JSON: %v", err)
	}

	// Проверка: JSON должен содержать ключи и значения из сделки и вложенных заказов
	jsonString := string(jsonData)
	expectedStrings := []string{
		`"id":1`,
		`"statusName":"completed"`,
		`"buyOrder":{"id":1,"userHashPublicKey":"buyer_public_key"`,
		`"sellOrder":{"id":2,"userHashPublicKey":"seller_public_key"`,
		`"createdAt":"2023-01-01T10:10:00Z"`,
		`"lastStatusChange":"2023-01-01T10:15:00Z"`,
	}
	for _, expected := range expectedStrings {
		if !strings.Contains(jsonString, expected) {
			t.Errorf("JSON representation is incorrect, missing: %s", expected)
		}
	}
}

func TestDealFromJson(t *testing.T) {
	// Подготовка JSON-данных
	jsonData := `{
		"id": 1,
		"buyOrder": {
			"id": 1,
			"userHashPublicKey": "buyer_public_key",
			"cryptocurrencyCode": "BTC",
			"typeName": "buy",
			"unitPrice": 50000.5,
			"quantity": 1.0,
			"createdAt": "2023-01-01T10:00:00Z",
			"lastStatusChange": "2023-01-01T10:05:00Z"
		},
		"sellOrder": {
			"id": 2,
			"userHashPublicKey": "seller_public_key",
			"cryptocurrencyCode": "BTC",
			"typeName": "sell",
			"unitPrice": 50000.5,
			"quantity": 1.0,
			"createdAt": "2023-01-01T09:50:00Z",
			"lastStatusChange": "2023-01-01T10:00:00Z"
		},
		"statusName": "completed",
		"createdAt": "2023-01-01T10:10:00Z",
		"lastStatusChange": "2023-01-01T10:15:00Z"
	}`

	// Тестирование метода FromJson
	deal, err := deal.FromJson([]byte(jsonData))
	if err != nil {
		t.Fatalf("Failed to parse JSON to deal: %v", err)
	}

	// Проверка полей объекта Deal
	if deal.ID != 1 {
		t.Errorf("Expected ID 1, got %d", deal.ID)
	}
	if deal.StatusName != "completed" {
		t.Errorf("Expected StatusName 'completed', got %s", deal.StatusName)
	}
	if deal.CreatedAt != "2023-01-01T10:10:00Z" {
		t.Errorf("Expected CreatedAt '2023-01-01T10:10:00Z', got %s", deal.CreatedAt)
	}
	if deal.LastStatusChange != "2023-01-01T10:15:00Z" {
		t.Errorf("Expected LastStatusChange '2023-01-01T10:15:00Z', got %s", deal.LastStatusChange)
	}

	// Проверка вложенных объектов Order
	if deal.BuyOrder == nil || deal.BuyOrder.ID != 1 {
		t.Errorf("Expected BuyOrder ID 1, got %v", deal.BuyOrder)
	}
	if deal.SellOrder == nil || deal.SellOrder.ID != 2 {
		t.Errorf("Expected SellOrder ID 2, got %v", deal.SellOrder)
	}
}
