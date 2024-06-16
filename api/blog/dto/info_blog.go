package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type InfoBlog struct {
	Title       string   `json:"title" validate:"required,min=3,max=500"`
	Description string   `json:"description" validate:"required,min=3,max=2000"`
	Slug        string   `json:"slug" validate:"required,min=3,max=200"`
	ImgURL      *string  `json:"imgUrl,omitempty" validate:"omitempty,uri,max=200"`
	Score       float64  `json:"score," validate:"required,min=0,max=1"`
	Tags        []string `json:"tags" validate:"required,dive,uppercase"`
}

func EmptyInfoBlog() *InfoBlog {
	return &InfoBlog{}
}

func (d *InfoBlog) GetValue() *InfoBlog {
	return d
}

func (b *InfoBlog) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	var msgs []string
	for _, err := range errs {
		switch err.Tag() {
		case "required":
			msgs = append(msgs, fmt.Sprintf("%s is required", err.Field()))
		case "min":
			msgs = append(msgs, fmt.Sprintf("%s must be at least %s characters", err.Field(), err.Param()))
		case "max":
			msgs = append(msgs, fmt.Sprintf("%s must be at most %s characters", err.Field(), err.Param()))
		case "url":
			msgs = append(msgs, fmt.Sprintf("%s must be a valid URL", err.Field()))
		case "uri":
			msgs = append(msgs, fmt.Sprintf("%s must be a valid URI", err.Field()))
		case "uppercase":
			msgs = append(msgs, fmt.Sprintf("%s must be uppercase", err.Field()))
		default:
			msgs = append(msgs, fmt.Sprintf("%s is invalid", err.Field()))
		}
	}
	return msgs, nil
}
