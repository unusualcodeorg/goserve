package network

import (
	"net/http"
)

type ResCode string

const (
	success_code              ResCode = "10000"
	failue_code               ResCode = "10001"
	retry_code                ResCode = "10002"
	invalid_access_token_code ResCode = "10003"
)

type responseModel struct {
	ResCode ResCode `json:"code" binding:"required"`
	Status  int     `json:"status" binding:"required"`
	Message string  `json:"message" binding:"required"`
	Data    any     `json:"data,omitempty" binding:"required,omitempty"`
}

func (r *responseModel) GetValue() *responseModel {
	return r
}

func SuccessDataResponse(message string, data any) ApiResponse {
	return &responseModel{
		ResCode: success_code,
		Status:  http.StatusOK,
		Message: message,
		Data:    data,
	}
}

func SuccessMsgResponse(message string) ApiResponse {
	return &responseModel{
		ResCode: success_code,
		Status:  http.StatusOK,
		Message: message,
	}
}

func BadRequestResponse(message string) ApiResponse {
	return &responseModel{
		ResCode: failue_code,
		Status:  http.StatusBadRequest,
		Message: message,
	}
}

func ForbiddenResponse(message string) ApiResponse {
	return &responseModel{
		ResCode: failue_code,
		Status:  http.StatusForbidden,
		Message: message,
	}
}

func UnauthorizedResponse(message string) ApiResponse {
	return &responseModel{
		ResCode: failue_code,
		Status:  http.StatusUnauthorized,
		Message: message,
	}
}

func NotFoundResponse(message string) ApiResponse {
	return &responseModel{
		ResCode: failue_code,
		Status:  http.StatusNotFound,
		Message: message,
	}
}

func InternalServerErrorResponse(message string) ApiResponse {
	return &responseModel{
		ResCode: failue_code,
		Status:  http.StatusInternalServerError,
		Message: message,
	}
}
