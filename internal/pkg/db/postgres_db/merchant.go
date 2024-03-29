package postgres_db

import (
	"fmt"
	"mockPay/internal/pkg/models"

	"github.com/jmoiron/sqlx"
)

type MerchantDB struct {
	db *sqlx.DB
}

func newMerchantDB(db *sqlx.DB) *MerchantDB {
	return &MerchantDB{
		db: db,
	}
}

func (r *MerchantDB) CraeteMerchant(merchant *models.Merchant) error {
	query := fmt.Sprintf("INSERT INTO %s (title, api_key) values ($1, $2) RETURNING id", merchantTable)

	row := r.db.QueryRow(query, merchant.Title, merchant.ApiKey)
	if err := row.Scan(&merchant.ID); err != nil {
		return err
	}
	return nil
}

func (r *MerchantDB) GetMerchant(merchant *models.Merchant) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", merchantTable)
	if err := r.db.Get(merchant, query, merchant.ID); err != nil {
		return err
	}
	return nil
}

func (r *MerchantDB) GetHashMerchant(merchant *models.Merchant) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE api_key=$1", merchantTable)
	if err := r.db.Get(merchant, query, merchant.ApiKey); err != nil {
		return err
	}
	return nil
}

func (r *MerchantDB) GetAllMerchant() (*[]models.Merchant, error) {
	merchants := []models.Merchant{}

	query := fmt.Sprintf("SELECT * FROM %s", merchantTable)
	if err := r.db.Select(&merchants, query); err != nil {
		return nil, err
	}
	return &merchants, nil
}

// TODO refactor
func (r *MerchantDB) MerchantTitle(title string) error {
	merchant := models.Merchant{}

	query := fmt.Sprintf("SELECT * FROM %s WHERE title=$1", merchantTable)
	if err := r.db.Get(&merchant, query, title); err != nil {
		return err
	}
	return nil
}

func (r *MerchantDB) CreateMerchantBalance(merchantBalance *models.MerchantBalance) error {
	query := fmt.Sprintf("INSERT INTO %s (merchant_id, merchant_balance) VALUES ($1, $2) RETURNING id",
		merchantBalanceTable)

	row := r.db.QueryRow(query, merchantBalance.MerchantID, merchantBalance.MerchantBalance)

	if err := row.Scan(&merchantBalance.ID); err != nil {
		return err
	}

	return nil
}
