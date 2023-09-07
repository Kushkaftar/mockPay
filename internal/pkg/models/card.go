package models

import "time"

type Card struct {
	ID         int       `db:"id" json:"id,omitempty"`
	PAN        int       `db:"pan" json:"pan" binding:"required"`
	CardHolder string    `db:"card_holder" json:"card_holder" binding:"required"`
	ExpMonth   int       `db:"exp_month" json:"exp_month" binding:"required"`
	ExpYear    int       `db:"exp_year" json:"exp_year" binding:"required"`
	CVC        int       `db:"cvc" json:"cvc" binding:"required"`
	HashCard   string    `db:"hash_card" json:"hash_card,omitempty"`
	CreatedAt  time.Time `db:"created_at"`
}
