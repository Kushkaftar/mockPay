package transaction_service

import (
	"errors"
	"log"
	"mockPay/internal/pkg/models"
	"mockPay/internal/services/balance_event"

	"github.com/google/uuid"
)

func (s *PurchaseService) Refund(refund *models.RefundRquest) (*models.RecurrentResponse, error) {

	// get card
	cardID, err := s.repository.GetCard(refund.CardHash, &refund.MerchantID)
	if err != nil {
		return nil, err
	}

	log.Printf("cardID - %d", cardID)

	// get transaction
	transactionToRefund := models.Transaction{
		MerchantID: refund.MerchantID,
		UUID:       refund.TransactonUUID,
	}

	if err := s.repository.Status(&transactionToRefund); err != nil {
		return nil, err
	}

	log.Printf("transactionToRefund - %+v", transactionToRefund)

	if transactionToRefund.CardID != cardID {
		return nil, errors.New("refund not possible")
	}

	if transactionToRefund.TransactionStatus == models.NewStatus ||
		transactionToRefund.TransactionStatus == models.ProcessingStatus {
		return nil, errors.New("target transaction not completed")
	}

	if transactionToRefund.TransactionStatus == models.RejectedStatus {
		return nil, errors.New("target transaction rejected")
	}

	if transactionToRefund.TransactionType != models.PurchaseType &&
		transactionToRefund.TransactionType != models.ReccurentType {
		log.Printf("target transaction status - %d", transactionToRefund.TransactionType)
		return nil, errors.New("target transaction type does not support refund")
	}

	// create transaction
	transaction := models.Transaction{
		MerchantID:        refund.MerchantID,
		CardID:            cardID,
		TransactionType:   models.RefundType,
		TransactionStatus: models.NewStatus,
		UUID:              uuid.New().String(),
		Amount:            refund.Amount,
	}

	// create refund struct
	refundDb := models.Refund{
		TargetTransactionID: transactionToRefund.ID,
	}

	// add to DB
	if err := s.repository.AddNewRefund(&transaction, &refundDb); err != nil {
		return nil, err
	}

	response := models.RecurrentResponse{
		Success:           true,
		TransactionType:   transactionType[transaction.TransactionType],
		TransactionStatus: transactionStatus[transaction.TransactionStatus],
		UUID:              transaction.UUID,
	}

	// TODO refactor
	// send postback
	go s.postback.SendPostback(transaction)

	// TODO del
	be := balance_event.NewBalanceEventService(s.allMethods, s.postback)

	go be.RefundBalanceEvent(&transaction, &transactionToRefund)

	return &response, nil
}
