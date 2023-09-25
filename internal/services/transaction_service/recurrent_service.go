package transaction_service

import (
	"mockPay/internal/pkg/models"
	"mockPay/internal/services/balance_event"

	"github.com/google/uuid"
)

func (s *PurchaseService) Recurrent(recurrent *models.Recurrent) (*models.RecurrentResponse, error) {

	// get card
	cardID, err := s.repository.GetCard(recurrent.HashCard, &recurrent.MerchantID)
	if err != nil {
		return nil, err
	}

	// create transaction
	transaction := models.Transaction{
		MerchantID:        recurrent.MerchantID,
		CardID:            cardID,
		TransactionType:   models.ReccurentType,
		TransactionStatus: models.NewStatus,
		UUID:              uuid.New().String(),
		Amount:            recurrent.Amount,
	}

	// add to DB
	if err := s.repository.AddNewRecurrent(&transaction); err != nil {
		return nil, err
	}

	recurrentResponse := models.RecurrentResponse{
		Success:           true,
		TransactionType:   transactionType[transaction.TransactionType],
		TransactionStatus: transactionStatus[transaction.TransactionStatus],
		UUID:              transaction.UUID,
	}

	// TODO refactor
	// send postback
	go s.postback.SendPostback(transaction)

	// TODO del
	bl := balance_event.NewBalanceEventService(s.allMethods, s.postback)

	go bl.PurchaseBalanceEvent(&transaction)

	return &recurrentResponse, nil
}
