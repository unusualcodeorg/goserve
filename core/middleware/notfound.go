package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

type notFoundMiddleware struct {
}

func NewNotFoundMiddleware() network.Middleware {
	m := notFoundMiddleware{}
	return &m
}

func (nfm *notFoundMiddleware) Mount(routeEngine *gin.Engine) {
	routeEngine.NoRoute(NotFoundHandler)
}

func NotFoundHandler(ctx *gin.Context) {
	network.NotFoundResponse("resource not found").Send(ctx)
}
