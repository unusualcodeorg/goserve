package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseCode string

const (
	success_code              ResponseCode = "10000"
	failue_code               ResponseCode = "10001"
	retry_code                ResponseCode = "10002"
	invalid_access_token_code ResponseCode = "10003"
)

type Response interface {
	Send(c *gin.Context)
}

type messageResponse struct {
	ResCode ResponseCode `json:"statusCode" binding:"required"`
	Status  int          `json:"status" binding:"required"`
	Message string       `json:"message" binding:"required"`
}

func (r messageResponse) Send(c *gin.Context) {
	c.JSON(int(r.Status), r)
}

type errorResponse struct {
	ResCode ResponseCode `json:"statusCode" binding:"required"`
	Status  int          `json:"status" binding:"required"`
	Message string       `json:"message" binding:"required"`
	Errors  []string     `json:"errors"`
}

func (r errorResponse) Send(c *gin.Context) {
	c.JSON(int(r.Status), r)
}

type dataResponse struct {
	ResCode ResponseCode `json:"statusCode" binding:"required"`
	Status  int          `json:"status" binding:"required"`
	Data    any          `json:"data" binding:"required"`
}

func (r dataResponse) Send(c *gin.Context) {
	c.JSON(int(r.Status), r)
}

func SuccessResponse(data any) Response {
	return &dataResponse{
		ResCode: success_code,
		Status:  http.StatusOK,
		Data:    data,
	}
}

func SuccessMsgResponse(message string) Response {
	return &messageResponse{
		ResCode: success_code,
		Status:  http.StatusOK,
		Message: message,
	}
}

func ErrorResponse(message string, status int, errors []string) Response {
	return &errorResponse{
		ResCode: failue_code,
		Status:  status,
		Message: message,
		Errors:  errors,
	}
}

func BadRequestResponse(errors []string) Response {
	return &errorResponse{
		ResCode: failue_code,
		Status:  http.StatusBadRequest,
		Message: "badRequest",
		Errors:  errors,
	}
}

func ForbiddenResponse(message string) Response {
	return &messageResponse{
		ResCode: failue_code,
		Status:  http.StatusForbidden,
		Message: message,
	}
}

func UnauthorizedResponse(message string) Response {
	return &messageResponse{
		ResCode: failue_code,
		Status:  http.StatusUnauthorized,
		Message: message,
	}
}

func NotFoundResponse(message string) Response {
	return &messageResponse{
		ResCode: failue_code,
		Status:  http.StatusNotFound,
		Message: message,
	}
}

func InternalServerErrorResponse(message string) Response {
	return &messageResponse{
		ResCode: failue_code,
		Status:  http.StatusInternalServerError,
		Message: message,
	}
}
