package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type SignUpBasic struct {
	Email         string `json:"email" binding:"required" validate:"required,email"`
	Password      string `json:"password" binding:"required" validate:"required,min=6,max=100"`
	Name          string `json:"name" binding:"required" validate:"required,min=2,max=200"`
	ProfilePicUrl *string `json:"profilePicUrl,omitempty" validate:"omitempty,url"`
}

func EmptySignUpBasic() *SignUpBasic {
	return &SignUpBasic{}
}

func (d *SignUpBasic) GetValue() *SignUpBasic {
	return d
}

func (d *SignUpBasic) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	var msgs []string
	for _, err := range errs {
		switch err.Tag() {
		case "required":
			msgs = append(msgs, fmt.Sprintf("%s is required", err.Field()))
		case "min":
			msgs = append(msgs, fmt.Sprintf("%s must be at least %s characters", err.Field(), err.Param()))
		case "max":
			msgs = append(msgs, fmt.Sprintf("%s must be at most %s characters", err.Field(), err.Param()))
		case "email":
			msgs = append(msgs, fmt.Sprintf("%s is not a valid email", err.Field()))
		case "url":
			msgs = append(msgs, fmt.Sprintf("%s is not a valid URL", err.Field()))
		default:
			msgs = append(msgs, fmt.Sprintf("%s is invalid", err.Field()))
		}
	}
	return msgs, nil
}
