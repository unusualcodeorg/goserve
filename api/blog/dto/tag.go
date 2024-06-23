package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func EmptyTag() *Tag {
	return &Tag{}
}

type Tag struct {
	Tag string `uri:"tag" validate:"required,uppercase"`
}

func (d *Tag) GetValue() *Tag {
	return d
}

func (b *Tag) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	var msgs []string
	for _, err := range errs {
		switch err.Tag() {
		case "required":
			msgs = append(msgs, fmt.Sprintf("%s is required", err.Field()))
		case "uppercase":
			msgs = append(msgs, fmt.Sprintf("%s must be uppercase", err.Field()))
		default:
			msgs = append(msgs, fmt.Sprintf("%s is invalid", err.Field()))
		}
	}
	return msgs, nil
}
