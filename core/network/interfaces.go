package network

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Response[T any] interface {
	GetValue() *T
	send(ctx *gin.Context)
}

type BaseController interface {
	Path() string
	Authentication() gin.HandlerFunc
	Authorization(role string) gin.HandlerFunc
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

type Param0MiddlewareProvider interface {
	Middleware() gin.HandlerFunc
}

type Param1MiddlewareProvider[T any] interface {
	Middleware(param1 T) gin.HandlerFunc
}

type Param2MiddlewareProvider[T any, V any] interface {
	Middleware(param1 T, param2 V) gin.HandlerFunc
}

type Param3MiddlewareProvider[T any, V any, W any] interface {
	Middleware(param1 T, param2 V, param3 W) gin.HandlerFunc
}

type ParamNMiddlewareProvider[T any] interface {
	Middleware(params ...T) gin.HandlerFunc
}

type AuthenticationProvider Param0MiddlewareProvider
type AuthorizationProvider ParamNMiddlewareProvider[string]

type Router interface {
	GetEngine() *gin.Engine
	RegisterValidationParsers(tagNameFunc validator.TagNameFunc)
	LoadControllers(controllers ...Controller)
	LoadRootMiddlewares(middlewares ...RootMiddleware)
	Start(ip string, port uint16)
}
