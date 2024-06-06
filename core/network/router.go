package network

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type router struct {
	routeEngine *gin.Engine
}

func NewRouter() Router {
	engine := gin.Default()
	r := router{
		routeEngine: engine,
	}
	return &r
}

func (r *router) GetRouteEngine() *gin.Engine {
	return r.routeEngine
}

func (r *router) LoadControllers(controllers ...Controller) {
	for _, c := range controllers {
		c.MountRoutes(r.routeEngine)
	}
}

func (r *router) LoadMiddlewares(middlewares ...Middleware) {
	for _, m := range middlewares {
		m.Mount(r.routeEngine)
	}
}

func (r *router) Start(ip string, port uint16) {
	address := fmt.Sprintf("%s:%d", ip, port)
	r.routeEngine.Run(address)
}
