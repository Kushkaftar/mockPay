package handlers

import (
	"log"
	"mockPay/internal/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) craeteMerchant(c *gin.Context) {
	merchant := models.Merchant{}

	if err := c.BindJSON(&merchant); err != nil {
		str := "json is not valid"
		newErrorResponse(c, http.StatusBadRequest, str)
		return
	}

	log.Printf("handler craeteMerchant, merchant - %+v", merchant)

	if err := h.merchantService.CraeteMerchant(&merchant); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("merchant - %+v", merchant)
	c.JSON(http.StatusOK, SuccessMerchant{Success: true, Merchant: merchant})
}

func (h *Handler) allMerchant(c *gin.Context) {
	merchants, err := h.merchantService.GetAllMerchant()
	if err != nil {
		log.Printf("handler allMerchant, err - %s", err)
	}
	log.Printf("handler allMerchant. merchants - %+v", merchants)
	c.JSON(http.StatusOK, SuccessMerchants{Success: true, Merchants: *merchants})
}
