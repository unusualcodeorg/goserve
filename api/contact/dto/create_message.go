package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type CreateMessage struct {
	Type string `json:"type" binding:"required,min=2,max=50"`
	Msg  string `json:"msg" binding:"required,min=0,max=2000"`
}

func EmptyCreateMessage() *CreateMessage {
	return &CreateMessage{}
}

func (d *CreateMessage) GetValue() *CreateMessage {
	return d
}

func (d *CreateMessage) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	var msgs []string
	for _, err := range errs {
		switch err.Tag() {
		case "required":
			msgs = append(msgs, fmt.Sprintf("%s is required", err.Field()))
		case "min":
			msgs = append(msgs, fmt.Sprintf("%s must be min %s", err.Field(), err.Param()))
		case "max":
			msgs = append(msgs, fmt.Sprintf("%s must be max%s", err.Field(), err.Param()))
		default:
			msgs = append(msgs, fmt.Sprintf("%s is invalid", err.Field()))
		}
	}
	return msgs, nil
}
