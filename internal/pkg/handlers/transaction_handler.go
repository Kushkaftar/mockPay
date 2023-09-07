package handlers

import (
	"log"
	"mockPay/internal/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) purchase(c *gin.Context) {

	purchase := models.PurchaseRequest{}

	if err := c.BindJSON(&purchase); err != nil {
		str := "json is not valid"
		newErrorResponse(c, http.StatusBadRequest, str)
		return
	}

	merchantID, err := getMerchantId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.transactionService.NewPurchase(purchase, merchantID)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) recurrent(c *gin.Context) {

	merchantID, err := getMerchantId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	recurrent := models.Recurrent{}
	if err := c.BindJSON(&recurrent); err != nil {
		str := "json is not valid"
		newErrorResponse(c, http.StatusBadRequest, str)
		return
	}

	recurrent.MerchantID = merchantID

	resp, err := h.transactionService.Recurrent(&recurrent)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) refund(c *gin.Context) {
	merchantID, err := getMerchantId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	refund := models.RefundRquest{
		MerchantID: merchantID,
	}

	if err := c.BindJSON(&refund); err != nil {
		str := "json is not valid"
		newErrorResponse(c, http.StatusBadRequest, str)
		return
	}

	resp, err := h.transactionService.Refund(&refund)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) status(c *gin.Context) {
	transaction := models.Transaction{}
	if err := c.BindJSON(&transaction); err != nil {
		str := "json is not valid"
		newErrorResponse(c, http.StatusBadRequest, str)
		return
	}

	merchantID, err := getMerchantId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	transaction.MerchantID = merchantID

	log.Printf("transaction - %+v", transaction)

	response, err := h.transactionService.Status(&transaction)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	log.Printf("after transaction - %+v", transaction)

	c.JSON(http.StatusOK, response)
}
