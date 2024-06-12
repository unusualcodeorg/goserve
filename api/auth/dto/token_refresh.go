package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type TokenRefresh struct {
	RefreshToken string `json:"refreshToken" binding:"required" validate:"required"`
}

func EmptyTokenRefresh() *TokenRefresh {
	return &TokenRefresh{}
}

func (d *TokenRefresh) GetValue() *TokenRefresh {
	return d
}

func (d *TokenRefresh) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	var msgs []string
	for _, err := range errs {
		switch err.Tag() {
		case "required":
			msgs = append(msgs, fmt.Sprintf("%s is required", err.Field()))
		default:
			msgs = append(msgs, fmt.Sprintf("%s is invalid", err.Field()))
		}
	}
	return msgs, nil
}
