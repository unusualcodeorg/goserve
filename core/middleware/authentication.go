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
	ctx.Next()
}
