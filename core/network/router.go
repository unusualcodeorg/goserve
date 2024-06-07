package network

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type router struct {
	engine *gin.Engine
}

func NewRouter(mode string) Router {
	gin.SetMode(mode)
	e := gin.Default()
	r := router{
		engine: e,
	}
	return &r
}

func (r *router) GetEngine() *gin.Engine {
	return r.engine
}

func (r *router) LoadControllers(controllers ...Controller) {
	for _, c := range controllers {
		g := r.engine.Group(c.Path())
		c.MountRoutes(g)
	}
}

func (r *router) LoadRootMiddlewares(middlewares ...RootMiddleware) {
	for _, m := range middlewares {
		m.Attach(r.engine)
	}
}

func (r *router) Start(ip string, port uint16) {
	address := fmt.Sprintf("%s:%d", ip, port)
	r.engine.Run(address)
}
