package parser

import (
	"github.com/go-playground/validator/v10"
)

func Validate(obj any) error {
	if err := validator.New().Struct(obj); err != nil {
		return err
	}
	return nil
}