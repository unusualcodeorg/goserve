package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/goserve/arch/network"
)

type notFound struct {
	network.BaseMiddleware
}

func NewNotFound() network.RootMiddleware {
	return &notFound{
		BaseMiddleware: network.NewBaseMiddleware(),
	}
}

func (m *notFound) Attach(engine *gin.Engine) {
	engine.NoRoute(m.Handler)
}

func (m *notFound) Handler(ctx *gin.Context) {
	m.Send(ctx).NotFoundError("resource not found", nil)
}
