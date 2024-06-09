package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/user/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/user/schema"
)

type UserAuth struct {
	User   *dto.InfoPrivateUser `json:"user" validate:"required"`
	Tokens *UserTokens          `json:"tokens" validate:"required"`
}

func NewUserAuthDto(user schema.User, tokens *UserTokens) *UserAuth {
	return &UserAuth{
		User:   dto.NewInfoPrivateUser(user),
		Tokens: tokens,
	}
}

func (d *UserAuth) Payload() *UserAuth {
	return d
}

func (d *UserAuth) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
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
