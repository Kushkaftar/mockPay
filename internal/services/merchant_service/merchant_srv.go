package merchant_service

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"mockPay/internal/pkg/db/postgres_db"
	"mockPay/internal/pkg/models"
	"time"
)

type MerchantSRV struct {
	repository *postgres_db.PostgresDB
}

func newMerchantSRV(repository *postgres_db.PostgresDB) *MerchantSRV {
	return &MerchantSRV{
		repository: repository,
	}
}

func (s *MerchantSRV) CraeteMerchant(merchant *models.Merchant) error {
	log.Printf("service CraeteMerchant, merchant - %+v", merchant)

	if err := s.repository.MerchantTitle(merchant.Title); err == nil {
		return errors.New("merchant title taken")
	} else if err != sql.ErrNoRows {
		return err
	}

	// create api key
	now := time.Now()
	strToHash := fmt.Sprintf("%s%d", merchant.Title, now.Unix())

	h := sha256.New()
	h.Write([]byte(strToHash))

	str := fmt.Sprintf("%x", h.Sum(nil))
	// end create api key

	merchant.ApiKey = str

	if err := s.repository.CraeteMerchant(merchant); err != nil {
		return err
	}

	// create merchant balance
	merchantBalance := models.MerchantBalance{
		MerchantID:      merchant.ID,
		MerchantBalance: 0,
	}

	if err := s.repository.CreateMerchantBalance(&merchantBalance); err != nil {
		return err
	}

	return nil
}

func (s *MerchantSRV) GetMerchant(merchant *models.Merchant) error {
	if err := s.repository.GetMerchant(merchant); err != nil {
		return err
	}

	return nil
}

func (s *MerchantSRV) GetAllMerchant() (*[]models.Merchant, error) {
	merchants, err := s.repository.GetAllMerchant()
	if err != nil {
		return nil, err
	}

	return merchants, nil
}

func (s *MerchantSRV) GetHashMerchant(merchant *models.Merchant) error {
	if err := s.repository.Merchant.GetHashMerchant(merchant); err != nil {
		return err
	}
	return nil
}
