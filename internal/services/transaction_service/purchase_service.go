package transaction_service

import (
	"errors"
	"fmt"
	"log"
	"mockPay/internal/pkg/db/postgres_db"
	"mockPay/internal/pkg/models"
	"mockPay/internal/services/balance_event"
	"mockPay/internal/services/postback"

	"github.com/google/uuid"
)

type PurchaseService struct {
	repository postgres_db.Transaction
	postback   *postback.Postback
	// TODO del allMethods
	allMethods *postgres_db.PostgresDB
}

func newPurchaseService(repository postgres_db.Transaction, allMethods *postgres_db.PostgresDB, postback *postback.Postback) *PurchaseService {
	return &PurchaseService{
		repository: repository,
		postback:   postback,
		allMethods: allMethods,
	}
}

func (s *PurchaseService) NewPurchase(purchase models.PurchaseRequest, merchantID int) (*models.PurchaseResponse, error) {

	// card model
	card := purchase.Card
	createCardHash(&card)

	// transaction model
	transaction := models.Transaction{
		MerchantID:        merchantID,
		TransactionType:   models.PurchaseType,
		TransactionStatus: models.NewStatus,
		UUID:              uuid.New().String(),
		Amount:            purchase.Amount,
	}

	// add to DB
	if err := s.repository.AddTransaction(&card, &transaction); err != nil {
		return nil, err
	}

	// create card balance
	cardBalance := models.CardBalance{
		CardID:      card.ID,
		CardBalance: getCardBalanceInPan(card.PAN),
	}

	if err := s.repository.CreateCardBalance(&cardBalance); err != nil {
		return nil, err
	}

	// TODO refactor
	purchaseResponse := models.PurchaseResponse{
		Success:           true,
		TransactionType:   transactionType[transaction.TransactionType],
		TransactionStatus: transactionStatus[transaction.TransactionStatus],
		UUID:              transaction.UUID,
		HashCard:          card.HashCard,
	}

	// TODO refactor
	// send postback
	go s.postback.SendPostback(transaction)

	// TODO del
	bl := balance_event.NewBalanceEventService(s.allMethods, s.postback)

	go bl.PurchaseBalanceEvent(&transaction)

	return &purchaseResponse, nil
}

func (s *PurchaseService) NewFormPurchase(purchase models.PurchaseFormRequest, merchantID int) (*models.PurchaseFormResponse, error) {
	// transaction model
	transaction := models.Transaction{
		MerchantID:        merchantID,
		TransactionType:   models.FormType,
		TransactionStatus: models.NewStatus,
		UUID:              uuid.New().String(),
		Amount:            purchase.Amount,
	}

	// add to DB
	if err := s.repository.FormTransaction(&transaction); err != nil {
		return nil, err
	}

	// TODO refactor
	formUrl := fmt.Sprintf("http://localhost:8080/api/form/%s", transaction.UUID)

	// TODO refactor
	purchaseFormResponse := models.PurchaseFormResponse{
		Success:           true,
		Amount:            transaction.Amount,
		TransactionType:   transactionType[transaction.TransactionType],
		TransactionStatus: transactionStatus[transaction.TransactionStatus],
		UUID:              transaction.UUID,
		Url:               formUrl,
	}

	// TODO refactor
	// send postback
	go s.postback.SendPostback(transaction)

	return &purchaseFormResponse, nil
}

// TODO refactor rename
func (s *PurchaseService) GetTransaction(transaction *models.Transaction) error {
	if err := s.repository.GetAmountTransaction(transaction); err != nil {
		return err
	}

	if transaction.TransactionStatus != models.NewStatus {
		return errors.New("operation is not supported for this transaction status")
	}

	return nil
}

func (s *PurchaseService) FormPurchase(card models.Card, transactoinUUID string) error {
	createCardHash(&card)

	tr := models.Transaction{
		UUID: transactoinUUID,
	}
	if err := s.repository.UpdateFormTransaction(&card, &tr); err != nil {
		return err
	}

	log.Printf("card - %+v", card)
	log.Printf("tr - %+v", tr)
	// create card balance
	cardBalance := models.CardBalance{
		CardID:      card.ID,
		CardBalance: getCardBalanceInPan(card.PAN),
	}

	if err := s.repository.CreateCardBalance(&cardBalance); err != nil {
		return err
	}

	// TODO refactor
	// send postback
	go s.postback.SendPostback(tr)

	// TODO del
	bl := balance_event.NewBalanceEventService(s.allMethods, s.postback)

	go bl.PurchaseBalanceEvent(&tr)

	return nil
}
