package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

type apikeyMiddleware struct {
}

func NewApikeyMiddleware() network.RootMiddleware {
	m := apikeyMiddleware{}
	return &m
}

func (m *apikeyMiddleware) Attach(engine *gin.Engine) {
	engine.Use(m.Handler)
}

func (m *apikeyMiddleware) Handler(ctx *gin.Context) {
	ctx.Next()
}
