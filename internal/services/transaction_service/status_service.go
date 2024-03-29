package transaction_service

import "mockPay/internal/pkg/models"

func (s *PurchaseService) Status(transaction *models.Transaction) (*models.PurchaseResponse, error) {

	if err := s.repository.Status(transaction); err != nil {
		return nil, err
	}

	resp := models.PurchaseResponse{
		Success:           true,
		TransactionType:   transactionType[transaction.TransactionType],
		TransactionStatus: transactionStatus[transaction.TransactionStatus],
		UUID:              transaction.UUID,
	}

	return &resp, nil
}
