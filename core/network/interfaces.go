package network

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BaseController interface {
	Path() string
	AuthenticationMiddleware() GroupMiddleware
	AuthorizationMiddleware() GroupMiddleware
}

type BaseService interface {
	Context() (context.Context, context.CancelFunc)
}

type Controller interface {
	BaseController
	MountRoutes(*gin.RouterGroup)
}

type Dto[T any] interface {
	Payload() *T
	ValidateErrors(validator.ValidationErrors) ([]string, error)
}

type RootMiddleware interface {
	Attach(*gin.Engine)
	Handler(*gin.Context)
}

type GroupMiddleware interface {
	Attach(*gin.RouterGroup)
	Handler(*gin.Context)
}

type GroupMiddlewareFunc func() GroupMiddleware

type Router interface {
	GetEngine() *gin.Engine
	LoadControllers(...Controller)
	LoadRootMiddlewares(...RootMiddleware)
	Start(ip string, port uint16)
}
