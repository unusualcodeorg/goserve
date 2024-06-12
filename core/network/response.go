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

func NewSuccessDataResponse(message string, data any) ApiResponse {
	return &responseModel{
		ResCode: success_code,
		Status:  http.StatusOK,
		Message: message,
		Data:    data,
	}
}

func NewSuccessMsgResponse(message string) ApiResponse {
	return &responseModel{
		ResCode: success_code,
		Status:  http.StatusOK,
		Message: message,
	}
}

func NewBadRequestResponse(message string) ApiResponse {
	return &responseModel{
		ResCode: failue_code,
		Status:  http.StatusBadRequest,
		Message: message,
	}
}

func NewForbiddenResponse(message string) ApiResponse {
	return &responseModel{
		ResCode: failue_code,
		Status:  http.StatusForbidden,
		Message: message,
	}
}

func NewUnauthorizedResponse(message string) ApiResponse {
	return &responseModel{
		ResCode: failue_code,
		Status:  http.StatusUnauthorized,
		Message: message,
	}
}

func NewNotFoundResponse(message string) ApiResponse {
	return &responseModel{
		ResCode: failue_code,
		Status:  http.StatusNotFound,
		Message: message,
	}
}

func NewInternalServerErrorResponse(message string) ApiResponse {
	return &responseModel{
		ResCode: failue_code,
		Status:  http.StatusInternalServerError,
		Message: message,
	}
}
