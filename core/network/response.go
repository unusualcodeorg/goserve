package network

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResCode string

const (
	success_code              ResCode = "10000"
	failue_code               ResCode = "10001"
	retry_code                ResCode = "10002"
	invalid_access_token_code ResCode = "10003"
)

type responseMessage struct {
	ResCode ResCode `json:"code" binding:"required"`
	Status  int     `json:"status" binding:"required"`
	Message string  `json:"message" binding:"required"`
}

func (r *responseMessage) GetValue() *responseMessage {
	return r
}

func (r *responseMessage) send(ctx *gin.Context) {
	ctx.JSON(int(r.Status), r)
}

type responseData struct {
	ResCode ResCode `json:"code" binding:"required"`
	Status  int     `json:"status" binding:"required"`
	Message string  `json:"message" binding:"required"`
	Data    any     `json:"data" binding:"required"`
}

func (r *responseData) GetValue() *responseData {
	return r
}

func (r *responseData) send(ctx *gin.Context) {
	ctx.JSON(int(r.Status), r)
}

func ResSuccessData(ctx *gin.Context, message string, data any) Response[responseData] {
	r := &responseData{
		ResCode: success_code,
		Status:  http.StatusOK,
		Message: message,
		Data:    data,
	}
	r.send(ctx)
	return r
}

func ResSuccessMsg(ctx *gin.Context, message string) Response[responseMessage] {
	r := &responseMessage{
		ResCode: success_code,
		Status:  http.StatusOK,
		Message: message,
	}
	r.send(ctx)
	return r
}

func ResBadRequest(ctx *gin.Context, message string) Response[responseMessage] {
	r := &responseMessage{
		ResCode: failue_code,
		Status:  http.StatusBadRequest,
		Message: message,
	}
	r.send(ctx)
	return r
}

func ResForbidden(ctx *gin.Context, message string) Response[responseMessage] {
	r := &responseMessage{
		ResCode: failue_code,
		Status:  http.StatusForbidden,
		Message: message,
	}
	r.send(ctx)
	return r
}

func ResUnauthorized(ctx *gin.Context, message string) Response[responseMessage] {
	r := &responseMessage{
		ResCode: failue_code,
		Status:  http.StatusUnauthorized,
		Message: message,
	}
	r.send(ctx)
	return r
}

func ResNotFound(ctx *gin.Context, message string) Response[responseMessage] {
	r := &responseMessage{
		ResCode: failue_code,
		Status:  http.StatusNotFound,
		Message: message,
	}
	r.send(ctx)
	return r
}

func ResInternalServerError(ctx *gin.Context, message string) Response[responseMessage] {
	r := &responseMessage{
		ResCode: failue_code,
		Status:  http.StatusInternalServerError,
		Message: message,
	}
	r.send(ctx)
	return r
}
