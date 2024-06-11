package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/auth"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

type authenticationProvider struct {
	authService auth.AuthService
}

func NewAuthenticationProvider(authService auth.AuthService) network.AuthenticationProvider {
	m := authenticationProvider{
		authService: authService,
	}
	return &m
}

func (m *authenticationProvider) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken := ctx.GetHeader(network.AuthorizationHeader)
		if len(accessToken) == 0 {
			panic(network.UnauthorizedError("permission denied: missing Authorization", nil))
		}
		ctx.Next()
	}
}
