package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/user/schema"
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
		_ = network.ReqMustGetUser[schema.User](ctx)

		ctx.Next()
	}
}
