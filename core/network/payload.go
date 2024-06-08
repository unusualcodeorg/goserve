package network

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
)

func Body(ctx *gin.Context, obj any) []string {
	if err := ctx.ShouldBindJSON(&obj); err != nil {
		errMsgs := getErrorMsgs(err)
		return errMsgs
	}
	return nil
}

func MapToDto(modelObj any, dtoObj any) any {
	copier.Copy(dtoObj, modelObj)
	return dtoObj
}

func getErrorMsgs(err error) []string {
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