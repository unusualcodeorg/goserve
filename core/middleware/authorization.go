package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

type authorizeProvider struct {
}

func NewAuthorizeProvider() network.MiddlewareProvider {
	m := authorizeProvider{}
	return &m
}

func (m *authorizeProvider) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}
