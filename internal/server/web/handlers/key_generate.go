package handlers

import (
	"encoding/json"
	"net/http"
	"sender/internal/data/blockchain/wallet"
)

type Keys struct {
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

func KeysGenerateHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	newWallet := wallet.New()
	walletSerialize := newWallet.Sereliaze()

	publicKey := walletSerialize.PublicKey
	privateKey := walletSerialize.PrivateKey

	jsonKeys := Keys{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jsonKeys)

}
