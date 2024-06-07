package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

type authorization struct {
}

func NewAuthorization() network.GroupMiddleware {
	m := authorization{}
	return &m
}

func (m *authorization) Attach(group *gin.RouterGroup) {
	group.Use(m.Handler)
}

func (m *authorization) Handler(ctx *gin.Context) {
	ctx.Next()
}
