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

	// get transaction
	transactionToRefund := models.Transaction{
		MerchantID: refund.MerchantID,
		UUID:       refund.TransactonUUID,
	}

	if err := s.repository.Status(&transactionToRefund); err != nil {
		return nil, err
	}

	if transactionToRefund.CardID != cardID {
		return nil, errors.New("refund not possible")
	}

	if transactionToRefund.TransactionStatus == newStatus ||
		transactionToRefund.TransactionStatus == processingStatus {
		return nil, errors.New("target transaction not completed")
	}

	if transactionToRefund.TransactionStatus == rejectedStatus {
		return nil, errors.New("target transaction rejected")
	}

	if transactionToRefund.TransactionType != purchaseType &&
		transactionToRefund.TransactionType != reccurentType {
		log.Printf("target transaction status - %d", transactionToRefund.TransactionType)
		return nil, errors.New("target transaction type does not support refund")
	}

	// create transaction
	transaction := models.Transaction{
		MerchantID:        refund.MerchantID,
		CardID:            cardID,
		TransactionType:   refundType,
		TransactionStatus: newStatus,
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

	// TODO del
	bl := balance_event.NewBalanceEventService(s.allMethods)

	go bl.RefundBalanceEvent(&transaction, &transactionToRefund)

	response := models.RecurrentResponse{
		Success:           true,
		TransactionType:   transactionType[transaction.TransactionType],
		TransactionStatus: transactionStatus[transaction.TransactionStatus],
		UUID:              transaction.UUID,
	}

	return &response, nil
}
