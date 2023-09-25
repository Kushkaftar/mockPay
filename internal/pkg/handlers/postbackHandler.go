package handlers

import (
	"log"
	"mockPay/internal/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) setPostback(c *gin.Context) {
	merchantID, err := getMerchantId(c)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	postback := models.Postback{
		MerchantID: merchantID,
	}

	if err := c.BindJSON(&postback); err != nil {
		str := "json is not valid"
		newErrorResponse(c, http.StatusBadRequest, str)
		return
	}

	log.Println(postback)

	if err := h.merchantService.SetPostback(&postback); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, SuccesPostback{Success: true})
}
