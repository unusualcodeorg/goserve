package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

type notFound struct {
}

func NewNotFound() network.RootMiddleware {
	m := notFound{}
	return &m
}

func (m *notFound) Attach(engine *gin.Engine) {
	engine.NoRoute(m.Handler)
}

func (*notFound) Handler(ctx *gin.Context) {
	network.ResNotFound(ctx, "resource not found")
}
