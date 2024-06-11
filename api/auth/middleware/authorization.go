package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/user/model"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

type authorizationProvider struct {
}

func NewAuthorizationProvider() network.AuthorizationProvider {
	m := authorizationProvider{}
	return &m
}

func (m *authorizationProvider) Middleware(roleNames ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(roleNames) == 0 {
			e := errors.New("permission denied: role missing")
			panic(network.ForbiddenError(e.Error(), e))
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
			e := errors.New("permission denied: does not have suffient role")
			panic(network.ForbiddenError(e.Error(), e))
		}

		ctx.Next()
	}
}
