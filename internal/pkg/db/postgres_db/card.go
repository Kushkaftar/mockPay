package postgres_db

import (
	"fmt"
	"log"
	"mockPay/internal/pkg/models"

	"github.com/jmoiron/sqlx"
)

type CardDB struct {
	db *sqlx.DB
}

func newCardDB(db *sqlx.DB) *CardDB {
	return &CardDB{
		db: db,
	}
}

func (r *CardDB) CraeteCard(card *models.Card) error {
	log.Printf("db CraeteMerchant, merchant - %+v", card)
	query := fmt.Sprintf("INSERT INTO %s (pan, card_holder, expMonth, expYear, cvc, hash_card) "+
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", cardTable)

	row := r.db.QueryRow(query,
		card.PAN,
		card.CardHolder,
		card.ExpMonth,
		card.ExpYear,
		card.CVC,
		card.HashCard)

	if err := row.Scan(&card.ID); err != nil {
		return err
	}
	return nil
}
