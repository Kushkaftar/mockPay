package models

type CardBalance struct {
	ID          int     `db:"id"`
	CardID      int     `db:"card_id"`
	CardBalance float32 `db:"card_balance"`
}

type MerchantBalance struct {
	ID              int     `db:"id"`
	MerchantID      int     `db:"merchant_id"`
	MerchantBalance float32 `db:"merchant_balance"`
}

type BalanceEvent struct {
	ID            int     `db:"id"`
	CustomerType  int     `db:"customer_type"`
	TransactionID int     `db:"transaction_id"`
	OldBalance    float32 `db:"old_balance"`
	NewBalance    float32 `db:"new_balance"`
}
