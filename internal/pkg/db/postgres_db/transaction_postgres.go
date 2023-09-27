package postgres_db

import (
	"database/sql"
	"fmt"
	"log"
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
		return err
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
		return err
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
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *TransactionDB) Status(transactoin *models.Transaction) error {
	var cardID sql.NullInt64

	query := fmt.Sprintf("SELECT id, transaction_type, transaction_status, card_id, amount FROM %s WHERE uuid=$1 AND merchant_id=$2;", transactionTable)
	row := r.db.QueryRow(query, transactoin.UUID, transactoin.MerchantID)
	if err := row.Scan(&transactoin.ID, &transactoin.TransactionType, &transactoin.TransactionStatus, &cardID, &transactoin.Amount); err != nil {
		return err
	}

	transactoin.CardID = int(cardID.Int64)

	return nil
}

func (r *TransactionDB) CreateCardBalance(cardBalance *models.CardBalance) error {
	query := fmt.Sprintf("INSERT INTO %s (card_id, card_balance) VALUES ($1, $2) RETURNING id",
		cardBalanceTable)

	row := r.db.QueryRow(query, cardBalance.CardID, cardBalance.CardBalance)

	if err := row.Scan(&cardBalance.ID); err != nil {
		return err
	}

	return nil
}

func (r *TransactionDB) UpdateTransactionStatus(transactoin *models.Transaction, status int) error {
	query := fmt.Sprintf("UPDATE %s SET status=$1 WHERE id=$2;", transactionTable)
	_, err := r.db.Exec(query, status, transactoin.ID)

	if err != nil {
		return err
	}

	transactoin.TransactionStatus = status

	return nil
}

func (r *TransactionDB) GetCard(cardHash string, merchantID *int) (int, error) {
	var id int
	query := fmt.Sprintf("select c.id  from %s c left join %s t on t.card_id = c.id where c.hash_card = $1 and t.merchant_id = $2 and t.transaction_type = 1", cardTable, transactionTable)
	row := r.db.QueryRow(query, cardHash, merchantID)

	if err := row.Scan(&id); err != nil {
		return 0, err
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
		return err
	}
	return nil
}

func (r *TransactionDB) AddNewRefund(transactoin *models.Transaction, refund *models.Refund) error {

	tx, err := r.db.Begin()
	if err != nil {
		return err
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
		return err
	}

	refund.RefundTransactionID = transactoin.ID
	log.Printf("transactoin - %+v", transactoin)
	log.Printf("refund - %+v", refund)

	createRefundQuery := fmt.Sprintf("INSERT INTO %s (target_transaction_id, refund_transaction_id) "+
		"VALUES ($1, $2) RETURNING id", refundTable)

	rowRefund := tx.QueryRow(createRefundQuery,
		refund.TargetTransactionID,
		refund.RefundTransactionID)
	err = rowRefund.Scan(&refund.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *TransactionDB) FormTransaction(transactoin *models.Transaction) error {
	query := fmt.Sprintf("INSERT INTO %s (merchant_id, card_id, transaction_type, transaction_status, uuid, amount) "+
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", transactionTable)

	row := r.db.QueryRow(query,
		transactoin.MerchantID,
		nil,
		transactoin.TransactionType,
		transactoin.TransactionStatus,
		transactoin.UUID,
		transactoin.Amount)

	if err := row.Scan(&transactoin.ID); err != nil {
		return err
	}

	return nil
}

// TODO refactor rename
func (r *TransactionDB) GetAmountTransaction(transactoin *models.Transaction) error {
	var amount sql.NullFloat64

	query := fmt.Sprintf("SELECT t.amount, t.transaction_status, t.transaction_type FROM %s t WHERE uuid=$1 and t.transaction_type = %d;", transactionTable, models.FormType)
	row := r.db.QueryRow(query, transactoin.UUID)
	if err := row.Scan(&amount, &transactoin.TransactionStatus, &transactoin.TransactionType); err != nil {
		return err
	}
	transactoin.Amount = float32(amount.Float64)
	return nil
}

func (r *TransactionDB) UpdateFormTransaction(card *models.Card, transactoin *models.Transaction) error {
	//begin transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// add card to db
	createCardQuery := fmt.Sprintf("INSERT INTO %s (pan, card_holder, exp_month, exp_year, cvc, hash_card) "+
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;", cardTable)

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
		return err
	}

	// add transaction to db
	transactoin.CardID = card.ID

	updateTransactionQuery := fmt.Sprintf("UPDATE %s SET card_id=$1 WHERE uuid=$2 RETURNING id, amount, merchant_id;", transactionTable)

	rowTransaction := tx.QueryRow(updateTransactionQuery, transactoin.CardID, transactoin.UUID)
	if err := rowTransaction.Scan(&transactoin.ID, &transactoin.Amount, &transactoin.MerchantID); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
