package balance_event

import (
	"log"
	"mockPay/internal/pkg/models"
)

func (s *BalanceEventService) RefundBalanceEvent(transaction *models.Transaction,
	targetTransaction *models.Transaction) {
	// 1 - проставляем транзакции статус processing
	// 1.1 - send postback
	// 2 - считаем все рефанды и если их сумма превышает сумму списания возвращаем ошибку
	// 3 - получаем баланс мерчанта, проверяем что на балансе достаточно средств
	// 3.1 - если средств не достаточно мнеяем статус на rejected, завершаем выполнение
	// 3.2 - send postback
	// 4 - если средств достаточно делаем запись в БД
	// 5 - устанавливаем статус complite
	// 5.1 - send postback

	// 1
	if err := s.repository.UpdateTransactionStatus(transaction, models.ProcessingStatus); err != nil {
		log.Printf("!!! ALARM !!! refundBalanceEvent, UpdateTransactionStatus err - %s", err)
		return
	}

	// 1.1
	go s.postback.SendPostback(*transaction)

	// 2
	// get all refand summ

	// get merchant balance
	merchantBalance := models.MerchantBalance{
		MerchantID: transaction.MerchantID,
	}

	if err := s.repository.GetMerchantBalance(&merchantBalance); err != nil {
		log.Printf("refundBalanceEvent, GetMerchantBalance, err - %s", err)
		return
	}

	refandSum, err := s.repository.GetSumAllRefands(targetTransaction.ID)
	if err != nil {
		log.Printf("refundBalanceEvent, GetSumAllRefands, err - %s", err)
		return
	}

	if *refandSum > targetTransaction.Amount {
		log.Print("refandSum > targetTransaction.Amount")

		if err := s.repository.UpdateTransactionStatus(transaction, models.RejectedStatus); err != nil {
			log.Printf("!!! ALARM !!! purchaseBalanceEvent, UpdateTransactionStatus err - %s", err)
			return
		}

		// 3.2
		go s.postback.SendPostback(*transaction)

		return
	}

	// get card balance
	cardBalance := models.CardBalance{
		CardID: transaction.CardID,
	}
	log.Printf("!!!!! purchaseBalanceEvent, cardBalance - %+v", cardBalance)
	if err := s.repository.GetCardBalance(&cardBalance); err != nil {
		log.Printf("purchaseBalanceEvent, GetCardBalance err - %s", err)
		return
	}

	// new merchant balance
	newMerchantBalance := merchantBalance.MerchantBalance - transaction.Amount
	merchantBalanceEvent := models.BalanceEvent{
		CustomerType:  merchantCustomer,
		TransactionID: transaction.ID,
		OldBalance:    merchantBalance.MerchantBalance,
		NewBalance:    newMerchantBalance,
	}

	newCardBalance := cardBalance.CardBalance - transaction.Amount
	cardBalanceEvent := models.BalanceEvent{
		CustomerType:  cardCustomer,
		TransactionID: transaction.ID,
		OldBalance:    cardBalance.CardBalance,
		NewBalance:    newCardBalance,
	}

	if err := s.repository.BalanceEvent(&merchantBalance,
		&merchantBalanceEvent,
		&cardBalance,
		&cardBalanceEvent); err != nil {
		log.Printf("purchaseBalanceEvent, BalanceEvent, err - %s", err)
		return
	}

	// 5
	if err := s.repository.UpdateTransactionStatus(transaction, models.ComplitedStatus); err != nil {
		log.Printf("!!! ALARM !!! purchaseBalanceEvent, UpdateTransactionStatus err - %s", err)
		return
	}

	// 5.1
	go s.postback.SendPostback(*transaction)

	return
}
