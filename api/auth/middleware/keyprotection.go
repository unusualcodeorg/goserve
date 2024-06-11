package middleware

import (

	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/auth"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

type keyProtection struct {
	authService auth.AuthService
}

func NewKeyProtection(authService auth.AuthService) network.RootMiddleware {
	m := keyProtection{
		authService: authService,
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

	apikey, err := m.authService.FindApiKey(key)
	if err != nil {
		panic(network.ForbiddenError("permission denied: invalid x-api-key", err))
	}

	network.ReqSetApiKey(ctx, apikey)

	ctx.Next()
}
