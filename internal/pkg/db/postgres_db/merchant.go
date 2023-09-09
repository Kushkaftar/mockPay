package postgres_db

import (
	"fmt"
	"log"
	e "mockPay/internal/pkg/errorWrap"
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
	log.Printf("db CraeteMerchant, merchant - %+v", merchant)
	query := fmt.Sprintf("INSERT INTO %s (title, api_key) values ($1, $2) RETURNING id", merchantTable)

	row := r.db.QueryRow(query, merchant.Title, merchant.ApiKey)
	if err := row.Scan(&merchant.ID); err != nil {
		return e.Wrap("merchantDB, CraeteMerchant", err)
	}
	return nil
}

func (r *MerchantDB) GetMerchant(merchant *models.Merchant) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", merchantTable)
	if err := r.db.Get(merchant, query, merchant.ID); err != nil {
		return e.Wrap("merchantDB, GetMerchant", err)
	}
	return nil
}

func (r *MerchantDB) GetHashMerchant(merchant *models.Merchant) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE api_key=$1", merchantTable)
	if err := r.db.Get(merchant, query, merchant.ApiKey); err != nil {
		return e.Wrap("merchantDB, GetHashMerchant", err)
	}
	return nil
}

func (r *MerchantDB) GetAllMerchant() (*[]models.Merchant, error) {
	merchants := []models.Merchant{}

	query := fmt.Sprintf("SELECT * FROM %s", merchantTable)
	if err := r.db.Select(&merchants, query); err != nil {
		return nil, e.Wrap("merchantDB, GetAllMerchant", err)
	}
	return &merchants, nil
}

func (r *MerchantDB) MerchantTitle(title string) error {
	merchant := models.Merchant{}

	query := fmt.Sprintf("SELECT * FROM %s WHERE title=$1", merchantTable)
	if err := r.db.Get(&merchant, query, title); err != nil {
		return e.Wrap("merchantDB, MerchantTitle", err)
	}
	return nil
}

func (r *MerchantDB) CreateMerchantBalance(merchantBalance *models.MerchantBalance) error {
	query := fmt.Sprintf("INSERT INTO %s (merchant_id, merchant_balance) VALUES ($1, $2) RETURNING id",
		merchantBalanceTable)

	row := r.db.QueryRow(query, merchantBalance.MerchantID, merchantBalance.MerchantBalance)

	if err := row.Scan(&merchantBalance.ID); err != nil {
		return e.Wrap("merchantDB, CreateMerchantBalance", err)
	}

	return nil
}
