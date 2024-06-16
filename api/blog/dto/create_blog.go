package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type CreateBlog struct {
	Title       string   `json:"title" validate:"required,min=3,max=500"`
	Description string   `json:"description" validate:"required,min=3,max=2000"`
	DraftText   string   `json:"draftText" validate:"required"`
	Text        string   `json:"text" validate:"required,max=50000"`
	Slug        string   `json:"slug" validate:"required,min=3,max=200"`
	ImgURL      string   `json:"imgUrl" validate:"required,uri,max=200"`
	Tags        []string `json:"tags" validate:"required,min=1,dive,uppercase"`
}

func EmptyCreateBlog() *CreateBlog {
	return &CreateBlog{}
}

func (d *CreateBlog) GetValue() *CreateBlog {
	return d
}

func (b *CreateBlog) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	var msgs []string
	for _, err := range errs {
		switch err.Tag() {
		case "required":
			msgs = append(msgs, fmt.Sprintf("%s is required", err.Field()))
		case "min":
			msgs = append(msgs, fmt.Sprintf("%s must be at least %s size", err.Field(), err.Param()))
		case "max":
			msgs = append(msgs, fmt.Sprintf("%s must be at most %s size", err.Field(), err.Param()))
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
