package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

type authenticationMiddleware struct {
}

func NewAuthenticationMiddleware() network.GroupMiddleware {
	m := authenticationMiddleware{}
	return &m
}

func (m *authenticationMiddleware) Attach(group *gin.RouterGroup) {
	group.Use(m.Handler)
}

func (m *authenticationMiddleware) Handler(ctx *gin.Context) {
	ctx.Next()
}
