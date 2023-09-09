package transaction_service

import (
	"mockPay/internal/pkg/db/postgres_db"
	"mockPay/internal/pkg/models"
)

var transactionType = map[int]string{
	models.PurchaseType:  "purchase",
	models.ReccurentType: "recurrent",
	models.RefundType:    "refund",
}

var transactionStatus = map[int]string{
	models.NewStatus:        "new",
	models.ProcessingStatus: "processing",
	models.ComplitedStatus:  "complite",
	models.RejectedStatus:   "rejected",
}

type Purchase interface {
	NewPurchase(purchase models.PurchaseRequest, merchantID int) (*models.PrchaseResponse, error)
	Recurrent(recurrent *models.Recurrent) (*models.RecurrentResponse, error)
	Refund(refund *models.RefundRquest) (*models.RecurrentResponse, error)
	Status(transaction *models.Transaction) (*models.PrchaseResponse, error)
}

type TransactionService struct {
	Purchase
}

func NewTransactionService(repository *postgres_db.PostgresDB) *TransactionService {
	return &TransactionService{
		Purchase: newPurchaseService(repository.Transaction, repository),
	}
}
