package postgres_db

import (
	"database/sql"
	"fmt"
	e "mockPay/internal/pkg/errorWrap"
	"mockPay/internal/pkg/models"

	"github.com/jmoiron/sqlx"
)

type TransactionDB struct {
	db *sqlx.DB
}

func newTransactionDB(db *sqlx.DB) *TransactionDB {
	return &TransactionDB{
		db: db,
	}
}

func (r *TransactionDB) AddTransaction(card *models.Card, transactoin *models.Transaction) error {
	tx, err := r.db.Begin()
	if err != nil {
		return e.Wrap("transactionDB, addTransaction, bd transactoin begin", err)
	}

	// add card to db
	createCardQuery := fmt.Sprintf("INSERT INTO %s (pan, card_holder, exp_month, exp_year, cvc, hash_card) "+
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", cardTable)

	rowCard := tx.QueryRow(createCardQuery,
		card.PAN,
		card.CardHolder,
		card.ExpMonth,
		card.ExpYear,
		card.CVC,
		card.HashCard)
	err = rowCard.Scan(&card.ID)
	if err != nil {
		tx.Rollback()
		return e.Wrap("transactionDB, addTransaction, insert card", err)
	}

	// add transaction to db
	transactoin.CardID = card.ID

	createTransactionQuery := fmt.Sprintf("INSERT INTO %s (merchant_id, card_id, transaction_type, transaction_status, uuid, amount) "+
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", transactionTable)

	rowTransaction := tx.QueryRow(createTransactionQuery,
		transactoin.MerchantID,
		transactoin.CardID,
		transactoin.TransactionType,
		transactoin.TransactionStatus,
		transactoin.UUID,
		transactoin.Amount)
	err = rowTransaction.Scan(&transactoin.ID)
	if err != nil {
		tx.Rollback()
		return e.Wrap("transactionDB, addTransaction, insert transactoin", err)
	}

	if err := tx.Commit(); err != nil {
		return e.Wrap("transactionDB, addTransaction, bd transactoin commit", err)
	}

	return nil
}

func (r *TransactionDB) Status(transactoin *models.Transaction) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE uuid=$1 AND merchant_id=$2;", transactionTable)
	if err := r.db.QueryRowx(query, transactoin.UUID, transactoin.MerchantID).StructScan(transactoin); err != nil {
		if err == sql.ErrNoRows {
			return e.ErrNotFound
		}
		return e.Wrap("transactionDB, status", err)
	}
	return nil
}

func (r *TransactionDB) CreateCardBalance(cardBalance *models.CardBalance) error {
	query := fmt.Sprintf("INSERT INTO %s (card_id, card_balance) VALUES ($1, $2) RETURNING id",
		cardBalanceTable)

	row := r.db.QueryRow(query, cardBalance.CardID, cardBalance.CardBalance)

	if err := row.Scan(&cardBalance.ID); err != nil {
		return e.Wrap("transactionDB, createCardBalance", err)
	}

	return nil
}

func (r *TransactionDB) UpdateTransactionStatus(transactoin *models.Transaction, status int) error {
	query := fmt.Sprintf("UPDATE %s SET status=$1 WHERE id=$2;", transactionTable)
	_, err := r.db.Exec(query, status, transactoin.ID)

	if err != nil {
		return e.Wrap("transactionDB, updateTransactionStatus", err)
	}

	transactoin.TransactionStatus = status

	return nil
}

func (r *TransactionDB) GetCard(cardHash string, merchantID *int) (int, error) {
	var id int
	query := fmt.Sprintf("select c.id  from %s c left join %s t on t.card_id = c.id where c.hash_card = $1 and t.merchant_id = $2 and t.transaction_type = 1", cardTable, transactionTable)
	row := r.db.QueryRow(query, cardHash, merchantID)

	if err := row.Scan(&id); err != nil {
		return 0, e.Wrap("transactionDB, getCard", err)
	}
	return id, nil
}

func (r *TransactionDB) AddNewRecurrent(transactoin *models.Transaction) error {
	createTransactionQuery := fmt.Sprintf("INSERT INTO %s (merchant_id, card_id, transaction_type, transaction_status, uuid, amount) "+
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", transactionTable)

	rowTransaction := r.db.QueryRow(createTransactionQuery,
		transactoin.MerchantID,
		transactoin.CardID,
		transactoin.TransactionType,
		transactoin.TransactionStatus,
		transactoin.UUID,
		transactoin.Amount)
	if err := rowTransaction.Scan(&transactoin.ID); err != nil {
		return e.Wrap("transactionDB, addNewRecurrent", err)
	}
	return nil
}

func (r *TransactionDB) AddNewRefund(transactoin *models.Transaction, refund *models.Refund) error {

	tx, err := r.db.Begin()
	if err != nil {
		return e.Wrap("transactionDB, addNewRefund, bd transactoin begin", err)
	}

	createTransactionQuery := fmt.Sprintf("INSERT INTO %s (merchant_id, card_id, transaction_type, transaction_status, uuid, amount) "+
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", transactionTable)

	rowTransaction := tx.QueryRow(createTransactionQuery,
		transactoin.MerchantID,
		transactoin.CardID,
		transactoin.TransactionType,
		transactoin.TransactionStatus,
		transactoin.UUID,
		transactoin.Amount)
	err = rowTransaction.Scan(&transactoin.ID)
	if err != nil {
		tx.Rollback()
		return e.Wrap("transactionDB, addNewRefund, create transactoin", err)
	}

	refund.RefundTransactionID = transactoin.ID

	createRefundQuery := fmt.Sprintf("INSERT INTO %s (target_transaction_id, refund_transaction_id) "+
		"VALUES ($1, $2) RETURNING id", refundTable)

	rowRefund := tx.QueryRow(createRefundQuery,
		refund.TargetTransactionID,
		refund.RefundTransactionID)
	err = rowRefund.Scan(&refund.ID)
	if err != nil {
		tx.Rollback()
		return e.Wrap("transactionDB, addNewRefund, crate refund", err)
	}

	if err := tx.Commit(); err != nil {
		return e.Wrap("transactionDB, addNewRefund, bd transactoin commit", err)
	}

	return nil
}
