package network

import "github.com/gin-gonic/gin"

type Controller interface {
	MountRoutes(router *gin.Engine)
}

type Router interface {
	GetRouteEngine() *gin.Engine
	LoadControllers(...Controller)
	LoadHandlers(notfound gin.HandlerFunc)
	Start(ip string, port uint16)
}
