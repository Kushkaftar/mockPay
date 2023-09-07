package app

import (
	"fmt"
	"mockPay/internal/pkg/db/postgres_db"
	"mockPay/internal/pkg/handlers"
	"mockPay/internal/pkg/models"
	"mockPay/internal/services/merchant_service"
	"mockPay/internal/services/transaction_service"
)

func Start(c *models.Config) error {

	fmt.Printf("config - %+v", c)

	connect, err := postgres_db.NewPostgresConnect(c.DB)
	if err != nil {
		return err
	}

	repositiry := postgres_db.NewPostgresDB(connect)

	merchantService := merchant_service.NewMerchantService(repositiry)
	transactionService := transaction_service.NewTransactionService(repositiry)

	handler := handlers.NewHandler(merchantService, transactionService)

	server := new(server)
	if err := server.serverRun(c.Server.Port, handler.InitRoutes()); err != nil {
		return err
	}

	return nil
}
