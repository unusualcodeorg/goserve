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
	authService auth.AuthService
	userService user.UserService
}

func NewAuthenticationProvider(authService auth.AuthService, userService user.UserService) network.AuthenticationProvider {
	m := authenticationProvider{
		authService: authService,
		userService: userService,
	}
	return &m
}

func (m *authenticationProvider) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(network.AuthorizationHeader)
		if len(authHeader) == 0 {
			panic(network.UnauthorizedError("permission denied: missing Authorization", nil))
		}

		token := utils.ExtractBearerToken(authHeader)
		if token == "" {
			panic(network.UnauthorizedError("permission denied: invalid Authorization", nil))
		}

		claims, err := m.authService.VerifyToken(token)
		if err != nil {
			panic(network.UnauthorizedError(err.Error(), err))
		}

		valid := m.authService.ValidateClaims(claims)
		if !valid {
			panic(network.UnauthorizedError("permission denied: invalid claims", nil))
		}

		userId, err := mongo.NewObjectID(claims.Subject)
		if err != nil {
			panic(network.UnauthorizedError("permission denied: invalid claims subject", nil))
		}

		user, err := m.userService.FindUserById(userId)
		if err != nil {
			panic(network.UnauthorizedError("permission denied: claims subject does not exists", nil))
		}

		keystore, err := m.authService.FindKeystore(user, claims.ID)
		if err != nil || keystore == nil {
			panic(network.UnauthorizedError("permission denied: invalid access token", nil))
		}

		network.ReqSetUser(ctx, user)

		ctx.Next()
	}
}
