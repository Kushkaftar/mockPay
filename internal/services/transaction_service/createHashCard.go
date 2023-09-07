package transaction_service

import (
	"crypto/md5"
	"fmt"
	"mockPay/internal/pkg/models"
	"time"
)

func createCardHash(card *models.Card) {
	now := time.Now()
	byteCard := []byte(fmt.Sprintf("%d %s %d %d %d %q",
		card.PAN,
		card.CardHolder,
		card.ExpMonth,
		card.ExpYear,
		card.CVC,
		now))

	hashCard := md5.Sum(byteCard)

	card.HashCard = fmt.Sprintf("%x", hashCard)
}
