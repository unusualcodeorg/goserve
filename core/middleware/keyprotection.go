package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

type keyProtection struct {
	secretService core.SecretService
}

func NewKeyProtection(secretService core.SecretService) network.RootMiddleware {
	m := keyProtection{
		secretService: secretService,
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

	apikey, err := m.secretService.FindApiKey(key)
	if err != nil {
		panic(network.ForbiddenError("permission denied: invalid x-api-key", err))
	}

	network.ReqSetApiKey(ctx, apikey)

	ctx.Next()
}
