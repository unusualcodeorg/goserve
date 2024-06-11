package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/user/model"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

type authorizationProvider struct {
}

func NewAuthorizationProvider() network.AuthorizationProvider {
	m := authorizationProvider{}
	return &m
}

func (m *authorizationProvider) Middleware(roleName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_ = network.ReqMustGetUser[model.User](ctx)

		ctx.Next()
	}
}
