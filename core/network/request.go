package network

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
)

const (
	reqPayloadApiKey string = "apikey"
	reqPayloadUser   string = "user"
)

func ReqSetApiKey[T any](ctx *gin.Context, value *T) {
	ctx.Set(reqPayloadApiKey, value)
}

func ReqMustGetApiKey[T any](ctx *gin.Context) *T {
	apikey, ok := ctx.MustGet(reqPayloadApiKey).(*T)
	if !ok {
		panic(errors.New("apikey missing for request"))
	}
	return apikey
}

func ReqSetUser[T any](ctx *gin.Context, value *T) {
	ctx.Set(reqPayloadUser, value)
}

func ReqMustGetUser[T any](ctx *gin.Context) *T {
	user, ok := ctx.MustGet(reqPayloadUser).(*T)
	if !ok {
		panic(errors.New("user missing for request"))
	}
	return user
}

// ShouldBindJSON in gin internally used go-playground/validator i.e. why we have error with validaiton info
func ReqBody[T any](ctx *gin.Context, dto Dto[T]) (*T, error) {
	if err := ctx.ShouldBindJSON(dto); err != nil {
		e := processErrors(dto, err)
		return nil, e
	}

	v := validator.New()
	v.RegisterTagNameFunc(CustomTagNameFunc())
	if err := v.Struct(dto); err != nil {
		e := processErrors(dto, err)
		return nil, e
	}

	return dto.GetValue(), nil
}

func ReqQuery[T any](ctx *gin.Context, dto Dto[T]) (*T, error) {
	if err := ctx.ShouldBindQuery(dto); err != nil {
		e := processErrors(dto, err)
		return nil, e
	}

	v := validator.New()
	v.RegisterTagNameFunc(CustomTagNameFunc())
	if err := v.Struct(dto); err != nil {
		e := processErrors(dto, err)
		return nil, e
	}

	return dto.GetValue(), nil
}

func ReqParams[T any](ctx *gin.Context, dto Dto[T]) (*T, error) {
	if err := ctx.ShouldBindUri(dto); err != nil {
		e := processErrors(dto, err)
		return nil, e
	}

	v := validator.New()
	v.RegisterTagNameFunc(CustomTagNameFunc())
	if err := v.Struct(dto); err != nil {
		e := processErrors(dto, err)
		return nil, e
	}

	return dto.GetValue(), nil
}

func ReqHeaders[T any](ctx *gin.Context, dto Dto[T]) (*T, error) {
	if err := ctx.ShouldBindHeader(dto); err != nil {
		e := processErrors(dto, err)
		return nil, e
	}

	v := validator.New()
	v.RegisterTagNameFunc(CustomTagNameFunc())
	if err := v.Struct(dto); err != nil {
		e := processErrors(dto, err)
		return nil, e
	}

	return dto.GetValue(), nil
}

func MapToDto[T any, V any](modelObj *V) (*T, error) {
	var dtoObj T
	err := copier.Copy(&dtoObj, modelObj)
	if err != nil {
		return nil, err
	}
	return &dtoObj, nil
}

func processErrors[T any](dto Dto[T], err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		msgs, e := dto.ValidateErrors(validationErrors)
		if e != nil {
			return e
		}
		var msg strings.Builder
		br := ", "
		for _, m := range msgs {
			msg.WriteString(m + br)
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
