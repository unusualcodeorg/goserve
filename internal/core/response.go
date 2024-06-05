package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatusCode string

const (
	STATUS_CODE_SUCCESS              StatusCode = "10000"
	STATUS_CODE_FAILURE              StatusCode = "10001"
	STATUS_CODE_RETRY                StatusCode = "10002"
	STATUS_CODE_INVALID_ACCESS_TOKEN StatusCode = "10003"
)

type ResponseStatus uint

const (
	STATUS_SUCCESS        ResponseStatus = http.StatusOK
	STATUS_BAD_REQUEST    ResponseStatus = http.StatusBadRequest
	STATUS_UNAUTHORIZED   ResponseStatus = http.StatusUnauthorized
	STATUS_FORBIDDEN      ResponseStatus = http.StatusForbidden
	STATUS_NOT_FOUND      ResponseStatus = http.StatusNotFound
	STATUS_INTERNAL_ERROR ResponseStatus = http.StatusInternalServerError
)

type MessageResponse struct {
	StatusCode StatusCode     `json:"statusCode" binding:"required"`
	Status     ResponseStatus `json:"status" binding:"required"`
	Message    string         `json:"message" binding:"required"`
}

func (r MessageResponse) Send(c *gin.Context) {
	c.JSON(int(r.Status), r)
}
