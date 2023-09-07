package models

import "time"

type Transaction struct {
	ID                int       `db:"id" json:"-"`
	MerchantID        int       `db:"merchant_id" json:"merchant_id"`
	CardID            int       `db:"card_id" json:"-"`
	TransactionType   int       `db:"transaction_type" json:"transaction_type"`
	TransactionStatus int       `db:"transaction_status" json:"transaction_status"`
	UUID              string    `db:"uuid" json:"uuid"`
	Amount            float32   `db:"amount" json:"amount"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
}
