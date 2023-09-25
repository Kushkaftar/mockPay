package handlers

import "mockPay/internal/pkg/models"

type SuccessMerchant struct {
	Success  bool            `json:"success"`
	Merchant models.Merchant `json:"merchant,omitempty"`
}

type SuccessMerchants struct {
	Success   bool              `json:"success"`
	Merchants []models.Merchant `json:"merchants,omitempty"`
}

type SuccesPostback struct {
	Success bool `json:"success"`
}
