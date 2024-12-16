package deal

import (
	"sender/data/order"
	"sender/internal/jsonutil"
)

type Deal struct {
	Id               int          `json:"id"`
	BuyOrder         *order.Order `json:"buyOrder"`
	SellOrder        *order.Order `json:"sellOrder"`
	StatusName       string       `json:"statusName"`
	CreatedAt        string       `json:"createdAt"`
	LastStatusChange string       `json:"lastStatusChange"`
}

func (d *Deal) ToJson() ([]byte, error) {
	return jsonutil.ToJSON(d)
}

func FromJson(json_b []byte) (*Deal, error) {
	var deal Deal
	err := jsonutil.FromJSON(json_b, &deal)

	return &deal, err
}

func (d *Deal) EncodeToString() string {
	return d.EncodeToString()
}
