package balance_event

import (
	"mockPay/internal/pkg/db/postgres_db"
	"mockPay/internal/services/postback"
)

const (
	_ = iota
	merchantCustomer
	cardCustomer
)

type BalanceEventService struct {
	repository postgres_db.Balance
	postback   *postback.Postback
}

func NewBalanceEventService(repository postgres_db.Balance, postback *postback.Postback) *BalanceEventService {
	return &BalanceEventService{
		repository: repository,
		postback:   postback,
	}
}
