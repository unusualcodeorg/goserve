package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

type keyProtection struct {
}

func NewKeyProtection() network.RootMiddleware {
	m := keyProtection{}
	return &m
}

func (m *keyProtection) Attach(engine *gin.Engine) {
	engine.Use(m.Handler)
}

func (m *keyProtection) Handler(ctx *gin.Context) {
	ctx.Next()
}
