package wallet

import (
	"sender/internal/jsonutil"
)

type WalletSerialize struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

func (ws *WalletSerialize) ToJson() ([]byte, error) {
	return jsonutil.ToJSON(ws)
}

func FromJson(json_b []byte) (*WalletSerialize, error) {
	var walletSerialize WalletSerialize

	err := jsonutil.FromJSON(json_b, &walletSerialize)
	return &walletSerialize, err
}
