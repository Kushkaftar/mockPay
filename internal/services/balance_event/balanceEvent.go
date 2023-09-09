package balance_event

import (
	"mockPay/internal/pkg/db/postgres_db"
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
