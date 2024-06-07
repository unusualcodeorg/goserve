package network

import "github.com/gin-gonic/gin"

type BaseController interface {
	Path() string
	AuthenticationMiddleware() GroupMiddleware
	AuthorizationMiddleware() GroupMiddleware
}

type Controller interface {
	BaseController
	MountRoutes(group *gin.RouterGroup)
}

type RootMiddleware interface {
	Attach(engine *gin.Engine)
	Handler(ctx *gin.Context)
}

type GroupMiddleware interface {
	Attach(group *gin.RouterGroup)
	Handler(ctx *gin.Context)
}

type GroupMiddlewareRecipe func() GroupMiddleware

type Router interface {
	GetEngine() *gin.Engine
	LoadControllers(...Controller)
	LoadRootMiddlewares(...RootMiddleware)
	Start(ip string, port uint16)
}
