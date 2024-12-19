package order

import (
	"sender/internal/jsonutil"
)

type Order struct {
	ID                 int     `json:"id"`
	UserLogin          string  `json:"userLogin"`
	WalletID           int     `json:"walletID"`
	CryptocurrencyCode string  `json:"cryptocurrencyCode"`
	CardID             int     `json:"cardID"`
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
