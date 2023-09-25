package handlers

import (
	"errors"
	"mockPay/internal/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	apiKeyHeader = "ApiKey"
	merchantCtx  = "merchantId"
)

func (h *Handler) authMerchant(c *gin.Context) {
	header := c.GetHeader(apiKeyHeader)

	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	merchant := models.Merchant{
		ApiKey: header,
	}

	if err := h.merchantService.Merchant.GetHashMerchant(&merchant); err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "unauth")
		return
	}

	// log.Printf("merchant - %+v", merchant)
	c.Set(merchantCtx, merchant.ID)
}

func getMerchantId(c *gin.Context) (int, error) {
	id, ok := c.Get(merchantCtx)
	if !ok {
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		return 0, errors.New("user id is of invalid type")
	}

	return idInt, nil
}
