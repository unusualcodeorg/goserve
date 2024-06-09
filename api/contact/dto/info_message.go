package dto

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type InfoMessage struct {
	ID        string    `json:"_id" binding:"required"`
	Type      string    `json:"type" binding:"required"`
	Msg       string    `json:"msg" binding:"required"`
	CreatedAt time.Time `json:"createdAt" binding:"required"`
}

func EmptyInfoMessage() *InfoMessage {
	return &InfoMessage{}
}

func (d *InfoMessage) Payload() *InfoMessage {
	return d
}

func (d *InfoMessage) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
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
