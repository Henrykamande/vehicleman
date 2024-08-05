package models

type CurrencySettings struct {
	ID           int    `json:"id"`
	UserID       int    `json:"user_id"`
	CurrencyCode string `json:"currency_code"`
}
