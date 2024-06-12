package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/auth"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

type keyProtection struct {
	network.ResponseSender
	authService auth.AuthService
}

func NewKeyProtection(authService auth.AuthService) network.RootMiddleware {
	return &keyProtection{
		ResponseSender: network.NewResponseSender(),
		authService:    authService,
	}
}

func (m *keyProtection) Attach(engine *gin.Engine) {
	engine.Use(m.Handler)
}

func (m *keyProtection) Handler(ctx *gin.Context) {
	key := ctx.GetHeader(network.ApiKeyHeader)
	if len(key) == 0 {
		m.Send(ctx).UnauthorizedError("permission denied: missing x-api-key header", nil)
		return
	}

	apikey, err := m.authService.FindApiKey(key)
	if err != nil {
		m.Send(ctx).ForbiddenError("permission denied: invalid x-api-key", err)
		return
	}

	network.ReqSetApiKey(ctx, apikey)

	ctx.Next()
}
