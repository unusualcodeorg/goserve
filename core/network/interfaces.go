package network

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Error[T any] interface {
	GetValue() *T
	Error() string
	Unwrap() error
}

type ApiError Error[apiError]

type Response[T any] interface {
	GetValue() *T
}

type ApiResponse Response[responseModel]

type ResponseSender[T any, V any] interface {
	Debug() bool
	SendResponse(ctx *gin.Context, response T)
	SendError(ctx *gin.Context, err V)
}

type ApiResponseSender ResponseSender[ApiResponse, ApiError]

type BaseController interface {
	ApiResponseSender
	Path() string
	Authentication() gin.HandlerFunc
	Authorization(role string) gin.HandlerFunc
}

type Controller interface {
	BaseController
	MountRoutes(group *gin.RouterGroup)
}

type BaseService interface {
	Context() (context.Context, context.CancelFunc)
}

type Dto[T any] interface {
	GetValue() *T
	ValidateErrors(errs validator.ValidationErrors) ([]string, error)
}

type BaseMiddleware interface {
	ApiResponseSender
	Debug() bool
}

type RootMiddleware interface {
	BaseMiddleware
	Attach(engine *gin.Engine)
	Handler(ctx *gin.Context)
}

type BaseMiddlewareProvider interface {
	BaseMiddleware
}

type Param0MiddlewareProvider interface {
	BaseMiddlewareProvider
	Middleware() gin.HandlerFunc
}

type Param1MiddlewareProvider[T any] interface {
	BaseMiddlewareProvider
	Middleware(param1 T) gin.HandlerFunc
}

type Param2MiddlewareProvider[T any, V any] interface {
	BaseMiddlewareProvider
	Middleware(param1 T, param2 V) gin.HandlerFunc
}

type Param3MiddlewareProvider[T any, V any, W any] interface {
	BaseMiddlewareProvider
	Middleware(param1 T, param2 V, param3 W) gin.HandlerFunc
}

type ParamNMiddlewareProvider[T any] interface {
	BaseMiddlewareProvider
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
