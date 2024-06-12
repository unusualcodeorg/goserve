package network

import (
	"errors"
	"fmt"
	"net/http"
)

type apiError struct {
	Code    int
	Message string
	Err     error
}

func (e *apiError) GetValue() *apiError {
	return e
}

func (e *apiError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%d - %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%d - %s", e.Code, e.Message)
}

func (e *apiError) Unwrap() error {
	return e.Err
}

func newApiError(code int, message string, err error) ApiError {
	apiError := apiError{
		Code:    code,
		Message: message,
		Err:     err,
	}
	if err == nil {
		apiError.Err = errors.New(message)
	}
	return &apiError
}

func BadRequestError(message string, err error) ApiError {
	return newApiError(http.StatusBadRequest, message, err)
}

func ForbiddenError(message string, err error) ApiError {
	return newApiError(http.StatusForbidden, message, err)
}

func UnauthorizedError(message string, err error) ApiError {
	return newApiError(http.StatusUnauthorized, message, err)
}

func NotFoundError(message string, err error) ApiError {
	return newApiError(http.StatusNotFound, message, err)
}

func InternalServerError(message string, err error) ApiError {
	return newApiError(http.StatusInternalServerError, message, err)
}
