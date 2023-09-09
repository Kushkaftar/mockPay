package balance_event

import (
	"log"
	"mockPay/internal/pkg/models"
)

func (s *BalanceEventService) PurchaseBalanceEvent(transaction *models.Transaction) {

	// 1 - проставляем транзакции статус processing
	// 2 - получаем баланс карты
	// 3 - проверяем что на балансе достаточно средст для списания
	// 3.1 - если средств не достаточно мнеяем статус на rejected, завершаем выполнение
	// 4 - если средств достаточно делаем запись в БД
	// 5 - устанавливаем статус complite

	// 1
	if err := s.repository.UpdateTransactionStatus(transaction, models.ProcessingStatus); err != nil {
		log.Printf("!!! ALARM !!! purchaseBalanceEvent, UpdateTransactionStatus err - %s", err)
		return
	}

	// 2
	cardBalance := models.CardBalance{
		CardID: transaction.CardID,
	}
	log.Printf("!!!!! purchaseBalanceEvent, cardBalance - %+v", cardBalance)
	if err := s.repository.GetCardBalance(&cardBalance); err != nil {
		log.Printf("purchaseBalanceEvent, GetCardBalance err - %s", err)
		return
	}

	// 3
	if cardBalance.CardBalance-transaction.Amount < 0 {

		// 3.1
		log.Println("there are not enough funds on the card")

		if err := s.repository.UpdateTransactionStatus(transaction, models.RejectedStatus); err != nil {
			log.Printf("!!! ALARM !!! purchaseBalanceEvent, UpdateTransactionStatus err - %s", err)
			return
		}

		return
	}

	// 4

	// get merchant balance
	merchantBalance := models.MerchantBalance{
		MerchantID: transaction.MerchantID,
	}

	if err := s.repository.GetMerchantBalance(&merchantBalance); err != nil {
		log.Printf("purchaseBalanceEvent, GetMerchantBalance, err - %s", err)
		return
	}

	newMerchantBalance := merchantBalance.MerchantBalance + transaction.Amount

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

	return
}
