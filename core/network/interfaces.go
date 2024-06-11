package network

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BaseController interface {
	Path() string
	Authentication() gin.HandlerFunc
	Authorization(roleCode string) gin.HandlerFunc
}

type BaseService interface {
	Context() (context.Context, context.CancelFunc)
}

type Controller interface {
	BaseController
	MountRoutes(group *gin.RouterGroup)
}

type Dto[T any] interface {
	GetValue() *T
	ValidateErrors(errs validator.ValidationErrors) ([]string, error)
}

type RootMiddleware interface {
	Attach(engine *gin.Engine)
	Handler(ctx *gin.Context)
}

type MiddlewareProvider interface {
	Middleware() gin.HandlerFunc
}

type ParamMiddlewareProvider[T any] interface {
	Middleware(param T) gin.HandlerFunc
}

type AuthenticationProvider MiddlewareProvider
type AuthorizationProvider ParamMiddlewareProvider[string]

type Router interface {
	GetEngine() *gin.Engine
	RegisterValidationParsers(tagNameFunc validator.TagNameFunc)
	LoadControllers(controllers ...Controller)
	LoadRootMiddlewares(middlewares ...RootMiddleware)
	Start(ip string, port uint16)
}
