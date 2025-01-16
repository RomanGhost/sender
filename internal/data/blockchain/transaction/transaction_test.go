package transaction_test

import (
	"sender/internal/data/blockchain/transaction"
	"sender/internal/data/blockchain/wallet"
	"sender/internal/data/deal"
	"sender/internal/data/order"
	"testing"
)

func TestTransactionCreation(t *testing.T) {
	// Создаем кошелек
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
	tx, err := transaction.New(w, d)
	if err != nil {
		t.Fatalf("Failed to create transaction: %v", err)
	}

	// Проверяем корректность данных транзакции
	if tx.Sender == "" || tx.BuyerPublicKey == "" || tx.SellerPublicKey == "" {
		t.Fatal("Transaction fields are not properly initialized")
	}
	if tx.Transfer != 5000.0 {
		t.Fatalf("Expected transfer amount to be 5000.0, got %v", tx.Transfer)
	}
}

func TestTransactionSigningAndVerification(t *testing.T) {
	// Создаем кошелек и транзакцию
	w := wallet.New()
	d := &deal.Deal{
		ID:        1,
		BuyOrder:  &order.Order{ID: 1, UserHashPublicKey: w.Sereliaze().PublicKey},
		SellOrder: &order.Order{ID: 2, UserHashPublicKey: w.Sereliaze().PublicKey},
	}
	tx, _ := transaction.New(w, d)

	// Подписываем транзакцию
	err := tx.Sign()
	if err != nil {
		t.Fatalf("Failed to sign transaction: %v", err)
	}

	// Проверяем подпись
	valid, err := tx.Verify(w.PublicKey)
	if err != nil {
		t.Fatalf("Failed to verify transaction: %v", err)
	}
	if !valid {
		t.Fatal("Transaction signature verification failed")
	}
}

func TestTransactionToJson(t *testing.T) {
	// Создаем кошелек и транзакцию
	w := wallet.New()
	d := &deal.Deal{
		ID:        1,
		BuyOrder:  &order.Order{ID: 1, UserHashPublicKey: w.Sereliaze().PublicKey},
		SellOrder: &order.Order{ID: 2, UserHashPublicKey: w.Sereliaze().PublicKey},
	}
	tx, _ := transaction.New(w, d)

	// Сериализуем транзакцию в JSON
	jsonData, err := tx.ToJson()
	if err != nil {
		t.Fatalf("Failed to serialize transaction to JSON: %v", err)
	}

	// Проверяем, что данные корректно сериализованы
	if len(jsonData) == 0 {
		t.Fatal("Serialized JSON is empty")
	}
}

func TestTransactionFromJson(t *testing.T) {
	// Создаем тестовые данные
	jsonData := `{
		"sender": "TestSender",
		"buyer": "TestBuyerPublicKey",
		"seller": "TestSellerPublicKey",
		"message": "{\"ID\":1,\"BuyOrder\":{\"ID\":1,\"UserHashPublicKey\":\"TestBuyerPublicKey\"},\"SellOrder\":{\"ID\":2,\"UserHashPublicKey\":\"TestSellerPublicKey\"}}",
		"transfer": 1000.0,
		"signature": "TestSignature"
	}`

	// Десериализуем транзакцию
	tx, err := transaction.FromJson([]byte(jsonData))
	if err != nil {
		t.Fatalf("Failed to deserialize transaction: %v", err)
	}

	// Проверяем корректность данных
	if tx.Sender != "TestSender" || tx.BuyerPublicKey != "TestBuyerPublicKey" {
		t.Fatal("Deserialized transaction contains incorrect data")
	}
	if tx.Transfer != 1000.0 {
		t.Fatalf("Expected transfer amount to be 1000.0, got %v", tx.Transfer)
	}
}

func TestInvalidTransactionSigning(t *testing.T) {
	var tx *transaction.Transaction

	// Попытка подписать nil транзакцию
	err := tx.Sign()
	if err == nil {
		t.Fatal("Expected error when signing nil transaction")
	}
}

func TestInvalidTransactionVerification(t *testing.T) {
	var tx *transaction.Transaction

	// Попытка проверить подпись nil транзакции
	_, err := tx.Verify(nil)
	if err == nil {
		t.Fatal("Expected error when verifying nil transaction")
	}
}
