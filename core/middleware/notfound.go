package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

type notFoundMiddleware struct {
}

func NewNotFoundMiddleware() network.RootMiddleware {
	m := notFoundMiddleware{}
	return &m
}

func (m *notFoundMiddleware) Attach(engine *gin.Engine) {
	engine.NoRoute(m.Handler)
}

func (*notFoundMiddleware) Handler(ctx *gin.Context) {
	network.NotFoundResponse("resource not found").Send(ctx)
}
