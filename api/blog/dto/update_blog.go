package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateBlog struct {
	ID          primitive.ObjectID `json:"_id" binding:"required" validate:"required"`
	Title       *string             `json:"title" validate:"omitempty,min=3,max=500"`
	Description *string             `json:"description" validate:"omitempty,min=3,max=2000"`
	DraftText   *string             `json:"draftText" validate:"omitempty,max=50000"`
	Slug        *string             `json:"slug" validate:"omitempty,min=3,max=200"`
	ImgURL      *string             `json:"imgUrl" validate:"omitempty,uri,max=200"`
	Tags        *[]string           `json:"tags" validate:"omitempty,min=1,dive,uppercase"`
}

func EmptyUpdateBlog() *UpdateBlog {
	return &UpdateBlog{}
}

func (d *UpdateBlog) GetValue() *UpdateBlog {
	return d
}

func (b *UpdateBlog) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
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
