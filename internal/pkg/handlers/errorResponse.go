package handlers

import "github.com/gin-gonic/gin"

type errorResponse struct {
	Error string `json:"error"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {

	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
