package app

import (
	"log"
	"mockPay/internal/pkg/db/postgres_db"
	"mockPay/internal/pkg/handlers"
	"mockPay/internal/pkg/models"
	"mockPay/internal/services/merchant_service"
	"mockPay/internal/services/postback"
	"mockPay/internal/services/transaction_service"
)

func MustStart(c *models.Config) {

	connect, err := postgres_db.NewPostgresConnect(c.DB)
	if err != nil {
		log.Fatalf("failed to connect to database, error - %s", err)
	}

	repository := postgres_db.NewPostgresDB(connect)

	postbackService := postback.NewPostback(repository.Postback)

	merchantService := merchant_service.NewMerchantService(repository)
	transactionService := transaction_service.NewTransactionService(repository, postbackService)

	handler := handlers.NewHandler(merchantService, transactionService)

	server := new(server)
	if err := server.serverRun(c.Server.Port, handler.InitRoutes()); err != nil {
		log.Fatalf("failed to raise the server, error - %s", err)
	}
}
