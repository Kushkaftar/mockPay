package models

type Recurrent struct {
	Amount     float32 `json:"amount" binding:"required"`
	HashCard   string  `json:"hash_card" binding:"required"`
	MerchantID int
}

type RecurrentResponse struct {
	Success           bool   `json:"success"`
	TransactionType   string `json:"transaction_type"`
	TransactionStatus string `json:"transaction_status"`
	UUID              string `json:"uuid"`
}
