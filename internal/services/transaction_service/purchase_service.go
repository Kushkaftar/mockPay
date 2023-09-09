package transaction_service

import (
	"mockPay/internal/pkg/db/postgres_db"
	"mockPay/internal/pkg/models"
	"mockPay/internal/services/balance_event"

	"github.com/google/uuid"
)

type PurchaseService struct {
	repository postgres_db.Transaction
	// TODO del allMethods
	allMethods *postgres_db.PostgresDB
}

func newPurchaseService(repository postgres_db.Transaction, allMethods *postgres_db.PostgresDB) *PurchaseService {
	return &PurchaseService{
		repository: repository,
		allMethods: allMethods,
	}
}

func (s *PurchaseService) NewPurchase(purchase models.PurchaseRequest, merchantID int) (*models.PrchaseResponse, error) {

	// card model
	card := purchase.Card
	createCardHash(&card)

	// transaction model
	transaction := models.Transaction{
		MerchantID:        merchantID,
		TransactionType:   models.PurchaseType,
		TransactionStatus: models.NewStatus,
		UUID:              uuid.New().String(),
		Amount:            purchase.Amount,
	}

	// add to DB
	if err := s.repository.AddTransaction(&card, &transaction); err != nil {
		return nil, err
	}

	// create card balance
	cardBalance := models.CardBalance{
		CardID:      card.ID,
		CardBalance: getCardBalanceInPan(card.PAN),
	}

	if err := s.repository.CreateCardBalance(&cardBalance); err != nil {
		return nil, err
	}

	// TODO del
	bl := balance_event.NewBalanceEventService(s.allMethods)

	go bl.PurchaseBalanceEvent(&transaction)

	// TODO refactor
	purchaseResponse := models.PrchaseResponse{
		Success:           true,
		TransactionType:   transactionType[transaction.TransactionType],
		TransactionStatus: transactionStatus[transaction.TransactionStatus],
		UUID:              transaction.UUID,
		HashCard:          card.HashCard,
	}
	return &purchaseResponse, nil
}
