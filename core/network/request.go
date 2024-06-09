package network

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
)

const (
	ReqPayloadApiKey string = "apikey"
	ReqPayloadUser   string = "user"
)

func ReqBody[T any](ctx *gin.Context) (*T, error) {
	var body T
	if err := ctx.ShouldBindJSON(&body); err != nil {
		e := parseError(err)
		return nil, e
	}
	return &body, nil
}

func ReqQuery[T any](ctx *gin.Context) (*T, error) {
	var query T
	if err := ctx.ShouldBindQuery(&query); err != nil {
		e := parseError(err)
		return nil, e
	}

	if err := validator.New().Struct(query); err != nil {
		e := parseError(err)
		return nil, e
	}

	return &query, nil
}

func ReqHeaders[T any](ctx *gin.Context) (*T, error) {
	var headers T
	if err := ctx.ShouldBindHeader(&headers); err != nil {
		e := parseError(err)
		return nil, e
	}

	if err := validator.New().Struct(headers); err != nil {
		e := parseError(err)
		return nil, e
	}

	return &headers, nil
}

func MapToDto[T any, V any](modelObj *V) (*T, error) {
	var dtoObj T
	err := copier.Copy(&dtoObj, modelObj)
	if err != nil {
		return nil, err
	}
	return &dtoObj, nil
}

func parseError(err error) error {
	var msg strings.Builder
	br := " | "

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, err := range validationErrors {
			switch err.Tag() {
			case "required":
				msg.WriteString(err.Field() + " is required" + br)
			case "min":
				msg.WriteString(err.Field() + " must be min " + err.Param() + br)
			case "max":
				msg.WriteString(err.Field() + " must be max " + err.Param() + br)
			default:
				msg.WriteString(err.Field() + " is invalid" + br)
			}
		}
		// Remove the trailing separator
		errorMsg := msg.String()
		if len(errorMsg) > 0 {
			errorMsg = errorMsg[:len(errorMsg)-len(br)]
		}
		return errors.New(errorMsg)
	}

	return err
}
