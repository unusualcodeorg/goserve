package network

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type router struct {
	engine *gin.Engine
}

func NewRouter() Router {
	engine := gin.Default()
	r := router{
		engine: engine,
	}
	return &r
}

func (r *router) GetRouteEngine() *gin.Engine {
	return r.engine
}

func (r *router) LoadControllers(controllers ...Controller) {
	for _, c := range controllers {
		c.MountRoutes(r.engine)
	}
}

func (r *router) LoadHandlers(notfound gin.HandlerFunc) {
	r.engine.NoRoute(notfound)
}

func (r *router) Start(ip string, port uint16) {
	address := fmt.Sprintf("%s:%d", ip, port)
	r.engine.Run(address)
}
