package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

type authentication struct {
}

func NewAuthentication() network.GroupMiddleware {
	m := authentication{}
	return &m
}

func (m *authentication) Attach(group *gin.RouterGroup) {
	group.Use(m.Handler)
}

func (m *authentication) Handler(ctx *gin.Context) {
	accessToken := ctx.GetHeader(network.AuthorizationHeader)
	if len(accessToken) == 0 {
		panic(network.UnauthorizedError("permission denied: missing Authorization", nil))
	}
	ctx.Next()
}
