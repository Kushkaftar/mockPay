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
		}
		payments := api.Group("payments", h.authApiMerchant)
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
