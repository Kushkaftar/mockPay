package handlers

import (
	"log"
	"mockPay/internal/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) formPagePayment(c *gin.Context) {
	parseUUID := c.Param("uuid")

	isUUID, err := uuid.Parse(parseUUID)
	if err != nil {
		log.Printf("parse uuid error - %s", err)
		c.HTML(http.StatusBadRequest, "404/index.tmpl", nil)
		return
	}

	transaction := models.Transaction{
		UUID: isUUID.String(),
	}

	if err := h.transactionService.GetTransaction(&transaction); err != nil {
		log.Printf("transaction not found, error - %s", err)
		c.HTML(http.StatusBadRequest, "404/index.tmpl", nil)
		return
	}

	c.HTML(http.StatusOK, "form/index.html", gin.H{"amount": transaction.Amount})
}

func (h *Handler) formPay(c *gin.Context) {
	parseUUID := c.Param("uuid")

	isUUID, err := uuid.Parse(parseUUID)
	if err != nil {
		log.Printf("parse uuid error - %s", err)
		c.HTML(http.StatusBadRequest, "404/index.tmpl", nil)
		return
	}

	transaction := models.Transaction{
		UUID: isUUID.String(),
	}

	if err := h.transactionService.GetTransaction(&transaction); err != nil {
		log.Printf("transaction not found, error - %s", err)
		c.HTML(http.StatusBadRequest, "404/index.tmpl", nil)
		return
	}

	formData := models.CardForm{}

	if err = c.Bind(&formData); err != nil {
		log.Printf("error - %s", err)
		c.HTML(http.StatusBadRequest, "404/index.tmpl", nil)
		return
	}

	card := models.Card{
		PAN:        formData.PAN,
		CardHolder: formData.CardHolder,
		ExpMonth:   formData.ExpMonth,
		ExpYear:    formData.ExpYear,
		CVC:        formData.CVC,
	}

	if err := h.transactionService.FormPurchase(card, isUUID.String()); err != nil {
		log.Printf("error - %s", err)
		c.HTML(http.StatusBadRequest, "fail/index.html", nil)
		return
	}

	log.Printf("card_numbaer- %+v", formData)
	log.Println("formPay")
	c.HTML(http.StatusOK, "success/index.html", nil)
}
