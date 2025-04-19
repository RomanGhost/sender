package block_test

import (
	"encoding/json"
	"sender/internal/data/blockchain/block"
	"sender/internal/data/blockchain/transaction"
	"testing"
	"time"
)

func generateTestTransaction(amount float64) transaction.Transaction {
	return transaction.Transaction{
		Sender:          "test_sender",
		BuyerPublicKey:  "test_buyer",
		SellerPublicKey: "test_seller",
		DealMessage:     "test_message",
		Transfer:        amount,
		Signature:       "test_signature",
	}
}

func TestBlock_ToJson(t *testing.T) {
	block := &block.Block{
		ID:           1,
		TimeCreated:  time.Now().UTC().Unix(),
		Transactions: []transaction.Transaction{generateTestTransaction(100)},
		PreviousHash: "abc123",
		Nonce:        42,
	}

	jsonData, err := block.ToJson()
	if err != nil {
		t.Fatalf("ToJson() returned unexpected error: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Unmarshaling JSON returned error: %v", err)
	}

	if result["id"] != float64(block.ID) {
		t.Errorf("Expected id %d, got %v", block.ID, result["id"])
	}
	if result["previous_hash"] != block.PreviousHash {
		t.Errorf("Expected previous_hash '%s', got '%s'", block.PreviousHash, result["previous_hash"])
	}
	if result["nonce"] != float64(block.Nonce) {
		t.Errorf("Expected nonce %d, got %v", block.Nonce, result["nonce"])
	}
}

func TestFromJSON(t *testing.T) {
	blockJson := []byte(`{
		"id": 1,
		"time_create": 1745089962,
		"transactions": [{"sender": "test_sender", "buyer": "test_buyer", "seller": "test_seller", "message": "test_message", "transfer": 100, "signature": "test_signature"}],
		"previous_hash": "abc123",
		"nonce": 42
	}`)

	block, err := block.FromJSON(blockJson)
	if err != nil {
		t.Fatalf("FromJSON() returned unexpected error: %v", err)
	}

	if block.ID != 1 {
		t.Errorf("Expected ID 1, got %d", block.ID)
	}
	if block.PreviousHash != "abc123" {
		t.Errorf("Expected PreviousHash 'abc123', got '%s'", block.PreviousHash)
	}
	if len(block.Transactions) != 1 {
		t.Fatalf("Expected 1 transaction, got %d", len(block.Transactions))
	}
	tx := block.Transactions[0]
	if tx.Sender != "test_sender" {
		t.Errorf("Expected sender 'test_sender', got '%s'", tx.Sender)
	}
	if tx.Transfer != 100 {
		t.Errorf("Expected transfer 100, got %f", tx.Transfer)
	}
}

func TestFromJSON_InvalidJSON(t *testing.T) {
	invalidJson := []byte(`{invalid}`)

	_, err := block.FromJSON(invalidJson)
	if err == nil {
		t.Error("Expected an error for invalid JSON, but got none")
	}
}

func TestBlock_GetTransactions(t *testing.T) {
	transactions := []transaction.Transaction{
		generateTestTransaction(100),
		generateTestTransaction(200),
	}

	block := &block.Block{
		Transactions: transactions,
	}

	result := block.GetTransactions()
	if len(result) != len(transactions) {
		t.Errorf("Expected %d transactions, got %d", len(transactions), len(result))
	}

	for i, tx := range result {
		if tx.Sender != transactions[i].Sender || tx.Transfer != transactions[i].Transfer {
			t.Errorf("Transaction mismatch at index %d: expected %+v, got %+v", i, transactions[i], tx)
		}
	}
}

func TestBlock_ToJsonEmptyBlock(t *testing.T) {
	block := &block.Block{}
	jsonData, err := block.ToJson()
	if err != nil {
		t.Fatalf("ToJson() returned unexpected error for empty block: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Unmarshaling JSON returned error: %v", err)
	}

	if result["id"] != float64(0) {
		t.Errorf("Expected id 0, got %v", result["id"])
	}
	if result["transactions"] == nil {
		t.Errorf("Expected transactions to be non-nil, got nil")
	}
}
