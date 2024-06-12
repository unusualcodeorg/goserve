package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
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
	m.SendResponse(ctx, network.NotFoundResponse("resource not found"))
}
