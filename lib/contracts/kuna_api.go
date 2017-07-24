package contracts

import "time"


type KunaOrder struct{
	ID int `json:"id"`
	Side string `json:"side"`
	OrdType string `json:"ord_type"`
	Price string `json:"price"`
	AvgPrice string `json:"avg_price"`
	State string `json:"state"`
	Market string `json:"market"`
	CreatedAt time.Time `json:"created_at"`
	Volume string `json:"volume"`
	RemainingVolume string `json:"remaining_volume"`
	ExpectedVolume string `json:"executed_volume"`
	TradesCount int `json:"trades_count"`
}

type KunaOrdersBook struct{
	Asks []KunaOrder `json:"asks"`
	Bids []KunaOrder `json:"bids"`
}