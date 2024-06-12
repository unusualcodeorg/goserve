package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/user/model"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

type authorizationProvider struct {
	network.ApiResponseSender
}

func NewAuthorizationProvider() network.AuthorizationProvider {
	return &authorizationProvider{
		ApiResponseSender: network.NewApiResponseSender(),
	}
}

func (m *authorizationProvider) Middleware(roleNames ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(roleNames) == 0 {
			m.SendError(ctx, network.ForbiddenError("permission denied: role missing", nil))
			return
		}

		user := network.ReqMustGetUser[model.User](ctx)

		hasRole := false
		for _, code := range roleNames {
			for _, role := range user.RoleDocs {
				if role.Code == model.RoleCode(code) {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			m.SendError(ctx, network.ForbiddenError("permission denied: does not have suffient role", nil))
			return
		}

		ctx.Next()
	}
}
