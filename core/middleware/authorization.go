package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/schema"
)

type authorizeProvider struct {
}

func NewAuthorizeProvider() network.AuthorizationProvider {
	m := authorizeProvider{}
	return &m
}

func (m *authorizeProvider) Middleware(roleCode schema.RoleCode) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_ = network.ReqGetUser(ctx)

		ctx.Next()
	}
}
