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

type messageResponse struct {
	ResCode ResCode `json:"code" binding:"required"`
	Status  int     `json:"status" binding:"required"`
	Message string  `json:"message" binding:"required"`
}

func (r *messageResponse) GetValue() *messageResponse {
	return r
}

func (r *messageResponse) send(ctx *gin.Context) {
	ctx.JSON(int(r.Status), r)
}

type dataResponse struct {
	ResCode ResCode `json:"code" binding:"required"`
	Status  int     `json:"status" binding:"required"`
	Message string  `json:"message" binding:"required"`
	Data    any     `json:"data" binding:"required"`
}

func (r *dataResponse) GetValue() *dataResponse {
	return r
}

func (r *dataResponse) send(ctx *gin.Context) {
	ctx.JSON(int(r.Status), r)
}

func SuccessDataResponse(ctx *gin.Context, message string, data any) Response[dataResponse] {
	r := &dataResponse{
		ResCode: success_code,
		Status:  http.StatusOK,
		Message: message,
		Data:    data,
	}
	r.send(ctx)
	return r
}

func SuccessMsgResponse(ctx *gin.Context, message string) Response[messageResponse] {
	r := &messageResponse{
		ResCode: success_code,
		Status:  http.StatusOK,
		Message: message,
	}
	r.send(ctx)
	return r
}

func BadRequestResponse(ctx *gin.Context, message string) Response[messageResponse] {
	r := &messageResponse{
		ResCode: failue_code,
		Status:  http.StatusBadRequest,
		Message: message,
	}
	r.send(ctx)
	return r
}

func ForbiddenResponse(ctx *gin.Context, message string) Response[messageResponse] {
	r := &messageResponse{
		ResCode: failue_code,
		Status:  http.StatusForbidden,
		Message: message,
	}
	r.send(ctx)
	return r
}

func UnauthorizedResponse(ctx *gin.Context, message string) Response[messageResponse] {
	r := &messageResponse{
		ResCode: failue_code,
		Status:  http.StatusUnauthorized,
		Message: message,
	}
	r.send(ctx)
	return r
}

func NotFoundResponse(ctx *gin.Context, message string) Response[messageResponse] {
	r := &messageResponse{
		ResCode: failue_code,
		Status:  http.StatusNotFound,
		Message: message,
	}
	r.send(ctx)
	return r
}

func InternalServerErrorResponse(ctx *gin.Context, message string) Response[messageResponse] {
	r := &messageResponse{
		ResCode: failue_code,
		Status:  http.StatusInternalServerError,
		Message: message,
	}
	r.send(ctx)
	return r
}
