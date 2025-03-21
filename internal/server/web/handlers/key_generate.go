package handlers

import (
	"net/http"
	"sender/internal/data/blockchain/wallet"

	"github.com/gin-gonic/gin"
)

type Keys struct {
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

func KeysGenerateHandler(c *gin.Context) {
	newWallet := wallet.New()
	walletSerialize := newWallet.Sereliaze()

	publicKey := walletSerialize.PublicKey
	privateKey := walletSerialize.PrivateKey

	jsonKeys := Keys{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}

	c.JSON(http.StatusOK, jsonKeys)
}
