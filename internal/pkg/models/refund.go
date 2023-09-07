package models

import "time"

type Refund struct {
	ID                  int       `db:"id"`
	TargetTransactionID int       `db:"target_transaction_id"`
	RefundTransactionID int       `db:"refund_transaction_id"`
	CreatedAt           time.Time `db:"created_at"`
}

type RefundRquest struct {
	CardHash       string  `json:"card_hash" binding:"required"`
	TransactonUUID string  `json:"transacton_uuid" binding:"required"`
	Amount         float32 `json:"amount" binding:"required"`
	MerchantID     int     `json:"-"`
}
