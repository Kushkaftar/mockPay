package postgres_db

import (
	"fmt"
	"mockPay/internal/pkg/models"

	"github.com/jmoiron/sqlx"
)

type PostbackDB struct {
	db *sqlx.DB
}

func newPostbackDB(db *sqlx.DB) *PostbackDB {
	return &PostbackDB{
		db: db,
	}
}

func (r *PostbackDB) CreatePostback(postback *models.Postback) error {
	query := fmt.Sprintf("INSERT INTO %s (merchant_id, postback_url, postback_method, is_enabled) values ($1, $2, $3, $4) RETURNING id", postbackTable)

	row := r.db.QueryRow(query, postback.MerchantID, postback.PostbackUrl, postback.PostbackMethod, postback.IsEnabled)
	if err := row.Scan(&postback.ID); err != nil {
		return err
	}
	return nil
}

func (r *PostbackDB) GetPostback(postback *models.Postback) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE merchant_id=$1", postbackTable)
	if err := r.db.Get(postback, query, postback.MerchantID); err != nil {
		return err
	}
	return nil
}

func (r *PostbackDB) UpdatePostback(postback *models.Postback) error {
	query := fmt.Sprintf("UPDATE %s SET postback_url=$1, postback_method=$2, is_enabled=$3 WHERE id=$4;", postbackTable)
	_, err := r.db.Exec(query, postback.PostbackUrl, postback.PostbackMethod, postback.IsEnabled, postback.ID)

	if err != nil {
		return err
	}

	return nil
}
