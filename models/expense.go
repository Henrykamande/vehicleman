package models

import "time"

type Expense struct {
	ExpenseID   int       `db:"expense_id" json:"expense_id"`
	VehicleID   int       `db:"vehicle_id" json:"vehicle_id"`
	CategoryID  int       `db:"category_id" json:"category_id"`
	Amount      float64   `db:"amount" json:"amount"`
	Description string    `db:"description,omitempty" json:"description,omitempty"`
	Receipt     string    `db:"receipt,omitempty" json:"receipt,omitempty"`
	ExpenseDate time.Time `db:"expense_date" json:"expense_date"`
}
