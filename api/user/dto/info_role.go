package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/user/schema"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InfoRole struct {
	ID   primitive.ObjectID `json:"_id" binding:"required" validate:"required"`
	Code schema.RoleCode    `json:"code" binding:"required" validate:"required,rolecode"`
}

func NewInfoRole(role schema.Role) *InfoRole {
	return &InfoRole{
		ID:   role.ID,
		Code: role.Code,
	}
}

func EmptyInfoRole() *InfoRole {
	return &InfoRole{}
}

func (d *InfoRole) Payload() *InfoRole {
	return d
}

func (d *InfoRole) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	var msgs []string
	for _, err := range errs {
		switch err.Tag() {
		case "required":
			msgs = append(msgs, fmt.Sprintf("%s is required", err.Field()))
		case "rolecode":
			msgs = append(msgs, fmt.Sprintf("%s missing %s", err.Field(), err.Param()))
		default:
			msgs = append(msgs, fmt.Sprintf("%s is invalid", err.Field()))
		}
	}
	return msgs, nil
}
