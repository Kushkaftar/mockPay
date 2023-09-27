package models

import "time"

type Card struct {
	ID         int       `db:"id" json:"id,omitempty"`
	PAN        int       `db:"pan" json:"pan" form:"pan" binding:"required"`
	CardHolder string    `db:"card_holder" json:"card_holder" form:"card_holder" binding:"required"`
	ExpMonth   int       `db:"exp_month" json:"exp_month" form:"exp_month" binding:"required"`
	ExpYear    int       `db:"exp_year" json:"exp_year" form:"exp_year" binding:"required"`
	CVC        int       `db:"cvc" json:"cvc" form:"cvc" binding:"required"`
	HashCard   string    `db:"hash_card" json:"hash_card,omitempty"`
	CreatedAt  time.Time `db:"created_at"`
}

type CardForm struct {
	PAN        int    `form:"pan"`
	CardHolder string `form:"card_holder"`
	ExpMonth   int    `form:"exp_month"`
	ExpYear    int    `form:"exp_year"`
	CVC        int    `form:"cvc"`
}
