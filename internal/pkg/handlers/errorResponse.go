package handlers

import "github.com/gin-gonic/gin"

type errorResponse struct {
	Success bool
	Errors  errorResp `json:"errors"`
}

type errorResp struct {
	Text string `json:"text"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {

	c.AbortWithStatusJSON(statusCode,
		errorResponse{Success: false, Errors: errorResp{Text: message}})
}
