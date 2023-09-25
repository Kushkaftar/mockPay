package handlers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(cors.Default())

	api := router.Group("/api")
	{
		merchant := api.Group("merchant")
		{
			merchant.GET("/all", h.allMerchant)
			merchant.POST("/create", h.craeteMerchant)

			postback := merchant.Group("postback", h.authMerchant)
			postback.POST("/set", h.setPostback)

			//merchant.POST("/set-postback", h.authMerchant, h.postbackMerchant)
		}
		payments := api.Group("payments", h.authMerchant)
		{
			transaction := payments.Group("/transaction")
			{
				transaction.POST("/purchase", h.purchase)
				transaction.POST("/recurrent", h.recurrent)
				transaction.POST("/refund", h.refund)
				transaction.POST("/status", h.status)
			}
		}
	}

	return router
}
