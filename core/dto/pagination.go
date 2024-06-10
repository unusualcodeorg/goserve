package coredto

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Pagination struct {
	Page  int64 `form:"page" binding:"required" validate:"required,min=1,max=1000"`
	Limit int64 `form:"limit" binding:"required" validate:"required,min=1,max=1000"`
}

func (d *Pagination) GetValue() *Pagination {
	return d
}

// strings.ToLower because gin query param validation does not give back form:"page"
func (d *Pagination) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	var msgs []string
	for _, err := range errs {
		switch err.Tag() {
		case "required":
			msgs = append(msgs, fmt.Sprintf("%s is required", strings.ToLower(err.Field())))
		case "min":
			msgs = append(msgs, fmt.Sprintf("%s must be min %s", strings.ToLower(err.Field()), err.Param()))
		case "max":
			msgs = append(msgs, fmt.Sprintf("%s must be max%s", strings.ToLower(err.Field()), err.Param()))
		default:
			msgs = append(msgs, fmt.Sprintf("%s is invalid", strings.ToLower(err.Field())))
		}
	}
	return msgs, nil
}
