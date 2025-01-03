package order

import (
	"sender/internal/jsonutil"
)

type Order struct {
	ID                 int     `json:"id"`
	UserHashPublicKey  string  `json:"userHashPublicKey"`
	CryptocurrencyCode string  `json:"cryptocurrencyCode"`
	TypeName           string  `json:"typeName"`
	UnitPrice          float64 `json:"unitPrice"`
	Quantity           float64 `json:"quantity"`
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
