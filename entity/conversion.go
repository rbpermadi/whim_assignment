package entity

import (
	"time"
)

//Currency data
type Conversion struct {
	ID             int64     `json:"id"`
	CurrencyIDFrom int64     `json:"currency_id_from"`
	CurrencyIDTo   int64     `json:"currency_id_to"`
	Rate           float64   `json:"rate"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
