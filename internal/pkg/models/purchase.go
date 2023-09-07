package models

type PurchaseRequest struct {
	Card   Card    `json:"card" binding:"required"`
	Amount float32 `json:"amount" binding:"required"`
}

type PrchaseResponse struct {
	Success           bool   `json:"success"`
	TransactionType   string `json:"transaction_type"`
	TransactionStatus string `json:"transaction_status"`
	UUID              string `json:"uuid"`
	HashCard          string `json:"hash_card,omitempty"`
}
