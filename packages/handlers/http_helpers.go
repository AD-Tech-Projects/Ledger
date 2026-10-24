package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func decodeRequestBody[T any](c *gin.Context) (*T, bool) {
	var value T

	if err := c.ShouldBindJSON(&value); err != nil {
		_ = c.Error(err)
		c.String(http.StatusBadRequest, "invalid request body")
		return nil, false
	}

	return &value, true
}

func encodeJSONResponse[T any](c *gin.Context, status int, result T) {
	c.JSON(status, result)
}

func writeError(c *gin.Context, status int, err error) {
	_ = c.Error(err)
	c.String(status, err.Error())
}

func writeInternalServerError(c *gin.Context, err error) {
	writeError(c, http.StatusInternalServerError, err)
}

func writeBadRequest(c *gin.Context, err error) {
	writeError(c, http.StatusBadRequest, err)
}

func writeNotFound(c *gin.Context, err error) {
	writeError(c, http.StatusNotFound, err)
}
