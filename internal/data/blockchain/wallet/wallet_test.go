package wallet_test

import (
	"bytes"
	"sender/internal/data/blockchain/wallet"
	"testing"
)

func TestWalletSerializationDeserialization(t *testing.T) {
	// Создаем новый кошелек
	w := wallet.New()

	// Сериализуем кошелек
	serialized := w.Sereliaze()

	// Десериализуем кошелек
	deserializedWallet, err := wallet.Deserialize(serialized)
	if err != nil {
		t.Fatalf("Failed to deserialize wallet: %v", err)
	}

	// Проверяем равенство публичных ключей
	if deserializedWallet.PublicKey.N.Cmp(w.PublicKey.N) != 0 {
		t.Fatal("Deserialized wallet public key does not match the original")
	}

	// Проверяем равенство приватных ключей
	if deserializedWallet.PrivateKey.D.Cmp(w.PrivateKey.D) != 0 {
		t.Fatal("Deserialized wallet private key does not match the original")
	}
}

func TestWalletSerializeToJson(t *testing.T) {
	// Создаем новый кошелек
	w := wallet.New()
	serialized := w.Sereliaze()

	// Преобразуем WalletSerialize в JSON
	jsonData, err := serialized.ToJson()
	if err != nil {
		t.Fatalf("Failed to serialize WalletSerialize to JSON: %v", err)
	}

	// Проверяем, что JSON содержит ожидаемые ключи
	if !bytes.Contains(jsonData, []byte(`"public_key"`)) || !bytes.Contains(jsonData, []byte(`"private_key"`)) {
		t.Fatal("Serialized JSON does not contain expected fields")
	}
}

func TestWalletSerializeFromJson(t *testing.T) {
	// Пример JSON для WalletSerialize
	jsonData := `{
		"public_key": "MHwwDQYJKoZIhvcNAQEBBQADawAwaAJhAP1kKiP5V+yjz7RtSzkZj6l+b6plfh8mFek0pr5/xwhBkDPWWFC8IHnE35YZ4EHv7zuyTgwnU1SRGGzQpMTsAeNj3SQIDAQAB",
		"private_key": "MIICWgIBAAKBgHlS9ySJPXhV6WmDoLBPqAIwPWTZ3QKi/NuP9Ccw/XwnBcTUIEExQFAZceFYOu6HFjWJzI+jTRCN0QtOZT8CF2GUS3cgWtBQpdiHxgTlGwJBgKwYgnsECwOBKJh8xXIvxmF0rxHLGEXZ3OGwERHXGkSt9qprPf8MeTTshyn5/OQxAgMBAAECgYA6wqugLhgQCpZZOGZskXaEmN/ZzJQbWz8fEZZW0DEJP0+8P4PPn4H6+Z+dDW8ldlMwvQTo/5gnr+JSfwiF+Tx69/Q5P4JSKLzQuLPkxEQX+/DNMyYkxZ+aYp8dj0V6syPtuGiOtBqJoOl8GStNVwtTpwtEAVtPA6CjHNk2HQJBAOaAjOMHyHUNFTw7VnPTrPs26GegDMFeorbc0JWdjWrwn+djtGn+oJ2r8kXeVcNZWHCiQoA0Vxhn5idHLxAwRWsCQQC+AcUJ1j6+RRfhPq5cu91sKzXdnvOOGb0OxwiVzYOTU1gwtp6H+dwfuFk4LkNMBQw/UrtfTIxfBgDbQGPB0wQbAkAgLjMa6Vw9fAxHB+erJHY8GGEtMyQL9MLUV5g+LD/bkRAib0UbCZc1tQoDfPl3xESONRL9m2jc5mt5AKCs9LBxAkBckxHZqoySgaQyBElAqCX4WYrNMQPSJwXmDQtEGDYiISMWJGL6OumZAcIryb2cL9OYBo9KLn4AOPI3hGgbx3WvAkB74I1RHhJeTJlLOemzON70qjYx9BNkp9sxE9z8kdpMQNB58kKAPnd2lRGng9aHdHIxJKl5QuRlAfqEX2oeTI4W"
	}`

	// Восстанавливаем WalletSerialize из JSON
	walletSerialize, err := wallet.FromJson([]byte(jsonData))
	if err != nil {
		t.Fatalf("Failed to deserialize WalletSerialize from JSON: %v", err)
	}

	// Проверяем, что ключи не пусты
	if walletSerialize.PublicKey == "" || walletSerialize.PrivateKey == "" {
		t.Fatal("Deserialized WalletSerialize contains empty keys")
	}
}
