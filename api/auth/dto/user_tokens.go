package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type UserTokens struct {
	AccessToken  string `json:"accessToken" binding:"required" validate:"required"`
	RefreshToken string `json:"refreshToken" binding:"required" validate:"required"`
}

func NewUserToken(access string, refresh string) *UserTokens {
	return &UserTokens{
		AccessToken:  access,
		RefreshToken: refresh,
	}
}

func (d *UserTokens) GetValue() *UserTokens {
	return d
}

func (d *UserTokens) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
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
