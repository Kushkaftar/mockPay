package merchant_service

import (
	"mockPay/internal/pkg/db/postgres_db"
	"mockPay/internal/pkg/models"
)

type Merchant interface {
	CraeteMerchant(merchant *models.Merchant) error
	GetMerchant(merchant *models.Merchant) error
	GetHashMerchant(merchant *models.Merchant) error
	GetAllMerchant() (*[]models.Merchant, error)
}

type MerchantService struct {
	Merchant
}

func NewMerchantService(repository *postgres_db.PostgresDB) *MerchantService {
	return &MerchantService{
		Merchant: newMerchantSRV(repository),
	}
}
