package network

import "github.com/gin-gonic/gin"

type Controller interface {
	MountRoutes(routeEngine *gin.Engine)
}

type Middleware interface {
	Mount(routeEngine *gin.Engine)
}

type Router interface {
	GetRouteEngine() *gin.Engine
	LoadControllers(...Controller)
	LoadMiddlewares(...Middleware)
	Start(ip string, port uint16)
}
