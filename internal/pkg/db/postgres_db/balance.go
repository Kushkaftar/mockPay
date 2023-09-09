package postgres_db

import (
	"fmt"
	"mockPay/internal/pkg/models"

	"github.com/jmoiron/sqlx"
)

type BalanceDB struct {
	db *sqlx.DB
}

func NewBalanceDB(db *sqlx.DB) *BalanceDB {
	return &BalanceDB{
		db: db,
	}
}

func (r *BalanceDB) GetCardBalance(cardBalance *models.CardBalance) error {
	cb := models.CardBalance{}

	query := fmt.Sprintf("SELECT * FROM %s WHERE card_id=$1;", cardBalanceTable)

	err := r.db.QueryRowx(query, cardBalance.CardID).StructScan(&cb)
	if err != nil {
		return err
	}

	cardBalance.CardBalance = cb.CardBalance
	cardBalance.ID = cb.ID
	return nil
}

func (r *BalanceDB) GetMerchantBalance(merchantBalance *models.MerchantBalance) error {
	//mb := models.MerchantBalance{}
	query := fmt.Sprintf("SELECT * FROM %s WHERE merchant_id=$1;", merchantBalanceTable)

	err := r.db.QueryRowx(query, merchantBalance.MerchantID).StructScan(merchantBalance)
	if err != nil {
		return err
	}

	return nil
}

func (r *BalanceDB) UpdateTransactionStatus(transactoin *models.Transaction, status int) error {
	query := fmt.Sprintf("UPDATE %s SET transaction_status=$1 WHERE id=$2;", transactionTable)
	_, err := r.db.Exec(query, status, transactoin.ID)

	if err != nil {
		return err
	}

	transactoin.TransactionStatus = status

	return nil
}

func (r *BalanceDB) BalanceEvent(
	merchantBalance *models.MerchantBalance,
	merchantBalanceEvent *models.BalanceEvent,
	cardBalance *models.CardBalance,
	cardBalanceEvent *models.BalanceEvent) error {

	// 1 - изменяем баланс карты
	// 2 - делаем запись в balance_event
	// 3 - изменяем баланс мерчанта
	// 4 - делаем запись в balance_event

	// transactoin BD start
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// 1
	cardBalanceQuery := fmt.Sprintf("UPDATE %s SET card_balance=$1 WHERE card_id=$2;", cardBalanceTable)
	_, err = tx.Exec(cardBalanceQuery, cardBalanceEvent.NewBalance, cardBalance.CardID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 2
	balanceEventQuery := fmt.Sprintf("INSERT INTO %s (customer_type, transaction_id, old_balance, new_balance) "+
		"VALUES ($1, $2, $3, $4) RETURNING id", balanceEventTable)
	rowCardEvent := tx.QueryRow(balanceEventQuery,
		cardBalanceEvent.CustomerType,
		cardBalanceEvent.TransactionID,
		cardBalanceEvent.OldBalance,
		cardBalanceEvent.NewBalance)

	err = rowCardEvent.Scan(&cardBalanceEvent.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 3
	merchantBalanceQuery := fmt.Sprintf("UPDATE %s SET merchant_balance=$1 WHERE merchant_id=$2;", merchantBalanceTable)
	_, err = tx.Exec(merchantBalanceQuery, merchantBalanceEvent.NewBalance, merchantBalance.MerchantID)
	if err != nil {
		tx.Rollback()
		return err
	}

	//4
	rowMerchantEvent := tx.QueryRow(balanceEventQuery,
		merchantBalanceEvent.CustomerType,
		merchantBalanceEvent.TransactionID,
		merchantBalanceEvent.OldBalance,
		merchantBalanceEvent.NewBalance)

	err = rowMerchantEvent.Scan(&merchantBalanceEvent.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	// end transactoin BD

	return nil
}

func (r *BalanceDB) GetSumAllRefands(targerTransactionID int) (*float32, error) {
	var sum float32
	query := fmt.Sprintf("select sum(t.amount)  from %s t where t.id in"+
		"(select r.refund_transaction_id from %s r where r.target_transaction_id= $1)"+
		"and t.transaction_status in (1,2,3)", transactionTable, refundTable)

	row := r.db.QueryRow(query, targerTransactionID)
	if err := row.Scan(&sum); err != nil {
		return nil, err
	}
	return &sum, nil
}
