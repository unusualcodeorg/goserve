package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/auth"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/user"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
	"github.com/unusualcodeorg/go-lang-backend-architecture/utils"
)

type authenticationProvider struct {
	network.ResponseSender
	authService auth.Service
	userService user.Service
}

func NewAuthenticationProvider(authService auth.Service, userService user.Service) network.AuthenticationProvider {
	return &authenticationProvider{
		ResponseSender: network.NewResponseSender(),
		authService:    authService,
		userService:    userService,
	}
}

func (m *authenticationProvider) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(network.AuthorizationHeader)
		if len(authHeader) == 0 {
			m.Send(ctx).UnauthorizedError("permission denied: missing Authorization", nil)
			return
		}

		token := utils.ExtractBearerToken(authHeader)
		if token == "" {
			m.Send(ctx).UnauthorizedError("permission denied: invalid Authorization", nil)
			return
		}

		claims, err := m.authService.VerifyToken(token)
		if err != nil {
			m.Send(ctx).UnauthorizedError(err.Error(), err)
			return
		}

		valid := m.authService.ValidateClaims(claims)
		if !valid {
			m.Send(ctx).UnauthorizedError("permission denied: invalid claims", nil)
			return
		}

		userId, err := mongo.NewObjectID(claims.Subject)
		if err != nil {
			m.Send(ctx).UnauthorizedError("permission denied: invalid claims subject", nil)
			return
		}

		user, err := m.userService.FindUserById(userId)
		if err != nil {
			m.Send(ctx).UnauthorizedError("permission denied: claims subject does not exists", err)
			return
		}

		keystore, err := m.authService.FindKeystore(user, claims.ID)
		if err != nil || keystore == nil {
			m.Send(ctx).UnauthorizedError("permission denied: invalid access token", err)
			return
		}

		network.ReqSetUser(ctx, user)
		network.ReqSetKeystore(ctx, keystore)

		ctx.Next()
	}
}
