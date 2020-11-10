package entity

//Currency data
type ConvertCurrencies struct {
	CurrencyIDFrom int64   `json:"currency_id_from"`
	CurrencyIDTo   int64   `json:"currency_id_to"`
	Amount         float64 `json:"amount"`
	Result         float64 `json:"result"`
}
