package models

import "time"

// Income represents the structure of the incomes table in the database
type Income struct {
	IncomeID    int       `json:"income_id"`
	VechileID   int       `json:"vehicle_id"`
	Amount      float64   `json:"amount"`
	PaymentDate time.Time `json:"payment_date"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
}
