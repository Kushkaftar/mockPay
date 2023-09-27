package transaction_service

import (
	"mockPay/internal/pkg/db/postgres_db"
	"mockPay/internal/pkg/models"
	"mockPay/internal/services/postback"
)

var transactionType = map[int]string{
	models.PurchaseType:  "purchase",
	models.ReccurentType: "recurrent",
	models.RefundType:    "refund",
	models.FormType:      "form",
}

var transactionStatus = map[int]string{
	models.NewStatus:        "new",
	models.ProcessingStatus: "processing",
	models.ComplitedStatus:  "complite",
	models.RejectedStatus:   "rejected",
}

type Purchase interface {
	NewPurchase(purchase models.PurchaseRequest, merchantID int) (*models.PurchaseResponse, error)
	Recurrent(recurrent *models.Recurrent) (*models.RecurrentResponse, error)
	Refund(refund *models.RefundRquest) (*models.RecurrentResponse, error)
	Status(transaction *models.Transaction) (*models.PurchaseResponse, error)
	NewFormPurchase(purchase models.PurchaseFormRequest, merchantID int) (*models.PurchaseFormResponse, error)
	GetTransaction(transaction *models.Transaction) error
	FormPurchase(card models.Card, transactoinUUID string) error
}

type TransactionService struct {
	Purchase
}

func NewTransactionService(repository *postgres_db.PostgresDB, postback *postback.Postback) *TransactionService {
	return &TransactionService{
		Purchase: newPurchaseService(repository.Transaction, repository, postback),
	}
}
