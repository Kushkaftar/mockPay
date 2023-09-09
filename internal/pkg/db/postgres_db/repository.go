package postgres_db

import (
	"mockPay/internal/pkg/models"

	"github.com/jmoiron/sqlx"
)

const (
	merchantTable        = "merchant"
	cardTable            = "card"
	transactionTable     = "transaction"
	refundTable          = "refund"
	cardBalanceTable     = "card_balance"
	merchantBalanceTable = "merchant_balance"
	balanceEventTable    = "balance_event"
)

type Merchant interface {
	CraeteMerchant(merchant *models.Merchant) error
	GetMerchant(merchant *models.Merchant) error
	GetHashMerchant(merchant *models.Merchant) error
	GetAllMerchant() (*[]models.Merchant, error)
	MerchantTitle(title string) error
	CreateMerchantBalance(merchantBalance *models.MerchantBalance) error
}

type Transaction interface {
	AddTransaction(card *models.Card, transactoin *models.Transaction) error
	Status(transactoin *models.Transaction) error
	CreateCardBalance(cardBalance *models.CardBalance) error
	GetCard(cardHash string, merchantID *int) (int, error)
	AddNewRecurrent(transactoin *models.Transaction) error
	AddNewRefund(transactoin *models.Transaction, refund *models.Refund) error
}

type Balance interface {
	GetCardBalance(cardBalance *models.CardBalance) error
	GetMerchantBalance(merchantBalance *models.MerchantBalance) error
	UpdateTransactionStatus(transactoin *models.Transaction, status int) error
	BalanceEvent(
		merchantBalance *models.MerchantBalance,
		merchantBalanceEvent *models.BalanceEvent,
		cardBalance *models.CardBalance,
		cardBalanceEvent *models.BalanceEvent) error
	GetSumAllRefands(targerTransactionID int) (*float32, error)
}

type PostgresDB struct {
	Merchant
	Transaction
	Balance
}

func NewPostgresDB(db *sqlx.DB) *PostgresDB {
	return &PostgresDB{
		Merchant:    newMerchantDB(db),
		Transaction: newTransactionDB(db),
		Balance:     NewBalanceDB(db),
	}
}
