package order

import (
	"sender/internal/jsonutil"
)

type Order struct {
	Id                 int     `json:"id"`
	UserLogin          string  `json:"userLogin"`
	WalletId           int     `json:"walletId"`
	CryptocurrencyCode string  `json:"cryptocurrencyCode"`
	CardId             int     `json:"cardId"`
	TypeName           string  `json:"typeName"`
	StatusName         string  `json:"statusName"`
	UnitPrice          float64 `json:"unitPrice"`
	Quantity           float64 `json:"quantity"`
	Description        string  `json:"description"`
	CreatedAt          string  `json:"createdAt"`
	LastStatusChange   string  `json:"lastStatusChange"`
}

func (o *Order) ToJson() ([]byte, error) {
	return jsonutil.ToJSON(o)
}

func FromJson(json_b []byte) (*Order, error) {
	var order Order

	err := jsonutil.FromJSON(json_b, &order)

	return &order, err
}
