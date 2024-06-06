package core

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Validate(obj any) error {
	if err := validator.New().Struct(obj); err != nil {
		return err
	}
	return nil
}

func ParseBody(ctx *gin.Context, obj any) []string {
	if err := ctx.ShouldBindJSON(&obj); err != nil {
		errMsgs := parseError(err)
		return errMsgs
	}
	return nil
}

func parseError(err error) []string {
	errMsgs := make([]string, 0)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, err := range validationErrors {
			switch err.Tag() {
			case "required":
				errMsgs = append(errMsgs, err.Field()+" is required")
			case "gt":
				errMsgs = append(errMsgs, err.Field()+" must be greater than "+err.Param())
			default:
				errMsgs = append(errMsgs, err.Field()+" is invalid")
			}
		}
		return errMsgs
	}
	return append(errMsgs, err.Error())
}
