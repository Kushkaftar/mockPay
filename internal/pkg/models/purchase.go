package models

type PurchaseRequest struct {
	Card   Card    `json:"card" binding:"required"`
	Amount float32 `json:"amount" binding:"required"`
}

type PurchaseResponse struct {
	Success           bool   `json:"success"`
	TransactionType   string `json:"transaction_type"`
	TransactionStatus string `json:"transaction_status"`
	UUID              string `json:"uuid"`
	HashCard          string `json:"hash_card,omitempty"`
}

type PurchaseFormRequest struct {
	Amount float32 `json:"amount" binding:"required"`
}

type PurchaseFormResponse struct {
	Amount            float32 `json:"amount" binding:"required"`
	Success           bool    `json:"success"`
	TransactionType   string  `json:"transaction_type"`
	TransactionStatus string  `json:"transaction_status"`
	UUID              string  `json:"uuid"`
	Url               string  `json:"url"`
}
