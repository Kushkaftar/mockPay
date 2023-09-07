package handlers

import (
	"mockPay/internal/services/merchant_service"
	"mockPay/internal/services/transaction_service"
)

type Handler struct {
	merchantService    *merchant_service.MerchantService
	transactionService *transaction_service.TransactionService
}

func NewHandler(
	service *merchant_service.MerchantService,
	transactionService *transaction_service.TransactionService) *Handler {
	return &Handler{
		merchantService:    service,
		transactionService: transactionService,
	}
}
