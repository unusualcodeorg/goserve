package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
	authModel "github.com/unusualcodeorg/go-lang-backend-architecture/api/auth/model"
	userModel "github.com/unusualcodeorg/go-lang-backend-architecture/api/user/model"
)

const (
	reqPayloadApiKey   string = "apikey"
	reqPayloadUser     string = "user"
	reqPayloadKeystore string = "keystore"
)

type GinContextUtil interface {
	SetApiKey(ctx *gin.Context, value *authModel.ApiKey)
	MustGetApiKey(ctx *gin.Context) *authModel.ApiKey
	SetUser(ctx *gin.Context, value *userModel.User)
	MustGetUser(ctx *gin.Context) *userModel.User
	SetKeystore(ctx *gin.Context, value *authModel.Keystore)
	MustGetKeystore(ctx *gin.Context) *authModel.Keystore
}

type ginContextUtil struct {}

func NewGinContextUtil() GinContextUtil {
	return &ginContextUtil{}
}

func (u *ginContextUtil) SetApiKey(ctx *gin.Context, value *authModel.ApiKey) {
	ctx.Set(reqPayloadApiKey, value)
}

func (u *ginContextUtil) MustGetApiKey(ctx *gin.Context) *authModel.ApiKey {
	value, ok := ctx.MustGet(reqPayloadApiKey).(*authModel.ApiKey)
	if !ok {
		panic(errors.New(reqPayloadApiKey + " missing for request"))
	}
	return value
}

func (u *ginContextUtil) SetUser(ctx *gin.Context, value *userModel.User) {
	ctx.Set(reqPayloadUser, value)
}

func (u *ginContextUtil) MustGetUser(ctx *gin.Context) *userModel.User {
	value, ok := ctx.MustGet(reqPayloadUser).(*userModel.User)
	if !ok {
		panic(errors.New(reqPayloadUser + " missing for request"))
	}
	return value
}

func (u *ginContextUtil) SetKeystore(ctx *gin.Context, value *authModel.Keystore) {
	ctx.Set(reqPayloadKeystore, value)
}

func (u *ginContextUtil) MustGetKeystore(ctx *gin.Context) *authModel.Keystore {
	value, ok := ctx.MustGet(reqPayloadKeystore).(*authModel.Keystore)
	if !ok {
		panic(errors.New(reqPayloadKeystore + " missing for request"))
	}
	return value
}
