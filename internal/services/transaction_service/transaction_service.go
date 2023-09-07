package transaction_service

import (
	"mockPay/internal/pkg/db/postgres_db"
	"mockPay/internal/pkg/models"
)

// transaction type
const (
	_ = iota
	purchaseType
	reccurentType
	refundType
)

var transactionType = map[int]string{
	purchaseType:  "purchase",
	reccurentType: "recurrent",
	refundType:    "refund",
}

// transaction stutus
const (
	_ = iota
	newStatus
	processingStatus
	complitedStatus
	rejectedStatus
)

var transactionStatus = map[int]string{
	newStatus:        "new",
	processingStatus: "processing",
	complitedStatus:  "complite",
	rejectedStatus:   "rejected",
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
