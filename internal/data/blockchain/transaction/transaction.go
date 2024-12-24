package transaction

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"sender/internal/data/blockchain/wallet"
	"sender/internal/data/deal"
	"sender/internal/jsonutil"
)

type Transaction struct {
	Sender      string  `json:"sender"`
	DealMessage string  `json:"message"` //*deal.Deal `json:"message"`//TODO(Решить что делать с полем сообщения)
	Transfer    float64 `json:"transfer"`
	Signature   string  `json:"signature"`
	wallet      *wallet.Wallet
	deal        *deal.Deal
}

// New creates a new transaction and initializes it with data
func New(walletKeys *wallet.Wallet, deal *deal.Deal) (*Transaction, error) {
	serializeWallet := walletKeys.Sereliaze()

	// Calculate transfer amount
	transfer := deal.BuyOrder.Quantity * deal.BuyOrder.UnitPrice
	jsonData, _ := deal.ToJson()
	dataString := string(jsonData)

	return &Transaction{
		Sender:      serializeWallet.PublicKey,
		DealMessage: dataString, //deal,
		Transfer:    transfer,
		Signature:   "",
		wallet:      walletKeys,
		deal:        deal,
	}, nil
}

// Sign signs the transaction using the private key
func (t *Transaction) Sign() error {
	if t == nil {
		return errors.New("transaction is nil")
	}

	// Format the data to sign (Sender, Message, Transfer)
	dataToSign := fmt.Sprintf("%s:%s:%v", t.Sender, t.DealMessage, t.Transfer)
	messageBytes := []byte(dataToSign)

	// Hash the data
	hasher := sha256.New()
	hasher.Write(messageBytes)
	hashedMessage := hasher.Sum(nil)

	// Create the signature
	privateKey := t.wallet.PrivateKey
	signatureBytes, err := rsa.SignPKCS1v15(rand.Reader, privateKey, 0, hashedMessage)
	if err != nil {
		return errors.New("failed to sign transaction: " + err.Error())
	}

	// Encode the signature in Base64
	t.Signature = base64.RawStdEncoding.EncodeToString(signatureBytes)
	return nil
}

// Verify verifies the signature of the transaction
func (t *Transaction) Verify(publicKey *rsa.PublicKey) (bool, error) {
	if t == nil {
		return false, errors.New("transaction is nil")
	}

	// Format the data to verify (Sender, Message, Transfer)
	dataToVerify := fmt.Sprintf("%s:%s:%v", t.Sender, t.DealMessage, t.Transfer)
	messageBytes := []byte(dataToVerify)

	// Hash the data
	hasher := sha256.New()
	hasher.Write(messageBytes)
	hashedMessage := hasher.Sum(nil)

	// Decode the signature from Base64
	signatureBytes, err := base64.RawStdEncoding.DecodeString(t.Signature)
	if err != nil {
		return false, errors.New("failed to decode signature: " + err.Error())
	}

	// Verify the signature
	err = rsa.VerifyPKCS1v15(publicKey, 0, hashedMessage, signatureBytes)
	if err != nil {
		return false, errors.New("signature verification failed: " + err.Error())
	}

	return true, nil
}

// ToJson serializes the transaction to JSON
func (t *Transaction) ToJson() ([]byte, error) {
	return jsonutil.ToJSON(t)
}

func (t *Transaction) GetDeal() *deal.Deal {
	return t.deal
}

// FromJson deserializes a transaction from JSON
func FromJson(jsonData []byte) (*Transaction, error) {
	var transaction Transaction
	err := jsonutil.FromJSON(jsonData, &transaction)
	if err != nil {
		return nil, err
	}

	// Вычленение сделки из поля сообщения
	dealJsonString := transaction.DealMessage
	dealFromJson, _ := deal.FromJson([]byte(dealJsonString))
	transaction.deal = dealFromJson

	return &transaction, nil
}
