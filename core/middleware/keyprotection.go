package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

type keyProtection struct {
	coreService core.CoreService
}

func NewKeyProtection(coreService core.CoreService) network.RootMiddleware {
	m := keyProtection{
		coreService: coreService,
	}
	return &m
}

func (m *keyProtection) Attach(engine *gin.Engine) {
	engine.Use(m.Handler)
}

func (m *keyProtection) Handler(ctx *gin.Context) {
	key := ctx.GetHeader(network.ApiKeyHeader)
	if len(key) == 0 {
		panic(network.UnauthorizedError("permission denied: missing x-api-key header", nil))
	}

	apikey, err := m.coreService.FindApiKey(key)
	if err != nil {
		panic(network.ForbiddenError("permission denied: invalid x-api-key", err))
	}

	ctx.Set(network.ReqPayloadApiKey, apikey)

	ctx.Next()
}
