package merchant_service

import (
	"database/sql"

	"mockPay/internal/pkg/models"
)

func (s *MerchantSRV) SetPostback(postback *models.Postback) error {

	check := models.Postback{
		MerchantID: postback.MerchantID,
	}

	// TODO: refactor
	err := s.repository.Postback.GetPostback(&check)

	if err != nil {
		if err == sql.ErrNoRows {
			if err := s.repository.Postback.CreatePostback(postback); err != nil {
				return err
			}

			return nil
		}

		return err
	}

	postback.ID = check.ID

	if err := s.repository.Postback.UpdatePostback(postback); err != nil {
		return err
	}

	return nil
}
