package wallet

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"log"
)

type Wallet struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

func New() *Wallet {
	// Генерация приватного ключа
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Не удалось сгенерировать приватный ключ: %v", err)
	}
	log.Println("Приватный ключ успешно сгенерирован.")

	// Получение публичного ключа из приватного
	publicKey := &privateKey.PublicKey
	log.Println("Публичный ключ успешно сгенерирован из приватного ключа.")

	return &Wallet{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
}

func (w *Wallet) Sereliaze() *WalletSerialize {
	publicKeyPKCS := x509.MarshalPKCS1PublicKey(w.PublicKey)
	privateKeyPKCS := x509.MarshalPKCS1PrivateKey(w.PrivateKey)

	publicKeyBase64 := base64.RawStdEncoding.EncodeToString(publicKeyPKCS)
	privateKeyBase64 := base64.RawStdEncoding.EncodeToString(privateKeyPKCS)

	return &WalletSerialize{
		PublicKey:  publicKeyBase64,
		PrivateKey: privateKeyBase64,
	}
}

func Deserialize(serialized *WalletSerialize) (*Wallet, error) {
	// Декодирование публичного ключа из base64
	publicKeyBytes, err := base64.RawStdEncoding.DecodeString(serialized.PublicKey)
	if err != nil {
		return nil, err
	}
	publicKey, err := x509.ParsePKCS1PublicKey(publicKeyBytes)
	if err != nil {
		return nil, err
	}

	// Декодирование приватного ключа из base64
	privateKeyBytes, err := base64.RawStdEncoding.DecodeString(serialized.PrivateKey)
	if err != nil {
		return nil, err
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBytes)
	if err != nil {
		return nil, err
	}

	// Восстановление объекта Wallet
	return &Wallet{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}, nil
}
