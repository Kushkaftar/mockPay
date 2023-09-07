package postgres_db

import (
	"fmt"
	"mockPay/internal/pkg/models"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func NewPostgresConnect(cfg models.DB) (*sqlx.DB, error) {
	dbConnect := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode)

	db, err := sqlx.Open("postgres", dbConnect)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
