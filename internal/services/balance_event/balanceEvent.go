package balance_event

import (
	"errors"
	"mockPay/internal/pkg/db/postgres_db"
	"mockPay/internal/pkg/models"
)

const (
	_ = iota
	purchaseType
	reccurentType
	refundType
)

const (
	_ = iota
	newStatus
	processingStatus
	complitedStatus
	rejectedStatus
)

const (
	_ = iota
	merchantCustomer
	cardCustomer
)

type BalanceEventService struct {
	repository postgres_db.Balance
}

func NewBalanceEventService(repository postgres_db.Balance) *BalanceEventService {
	return &BalanceEventService{
		repository: repository,
	}
}

func (s *BalanceEventService) BalanceEvent(transaction *models.Transaction) error {
	switch transaction.TransactionType {
	// purchase
	case purchaseType:
		if err := s.purchaseBalanceEvent(transaction); err != nil {
			return err
		}

	// reccurent
	case reccurentType:
		if err := s.purchaseBalanceEvent(transaction); err != nil {
			return err
		}

	}

	return errors.New("out of range")
}
