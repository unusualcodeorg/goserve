package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

type authorizationMiddleware struct {
}

func NewAuthorizationMiddleware() network.GroupMiddleware {
	m := authorizationMiddleware{}
	return &m
}

func (m *authorizationMiddleware) Attach(group *gin.RouterGroup) {
	group.Use(m.Handler)
}

func (m *authorizationMiddleware) Handler(ctx *gin.Context) {
	ctx.Next()
}
