package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

type authProvider struct {
}

func NewAuthProvider() network.AuthenticationProvider {
	m := authProvider{}
	return &m
}

func (m *authProvider) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken := ctx.GetHeader(network.AuthorizationHeader)
		if len(accessToken) == 0 {
			panic(network.UnauthorizedError("permission denied: missing Authorization", nil))
		}
		ctx.Next()
	}
}
