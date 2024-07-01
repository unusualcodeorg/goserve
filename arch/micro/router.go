package micro

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/unusualcodeorg/goserve/arch/network"
)

type router struct {
	netRouter  network.Router
	natsClient *NatsClient
}

func NewRouter(mode string, natsClient *NatsClient) Router {
	return &router{
		netRouter:  network.NewRouter(mode),
		natsClient: natsClient,
	}
}

func (r *router) GetEngine() *gin.Engine {
	return r.netRouter.GetEngine()
}

func (r *router) NatsClient() *NatsClient {
	return r.natsClient
}

func (r *router) LoadRootMiddlewares(middlewares []network.RootMiddleware) {
	r.netRouter.LoadRootMiddlewares(middlewares)
}

func (r *router) LoadControllers(controllers []Controller) {
	nc := make([]network.Controller, len(controllers))
	for i, c := range controllers {
		nc[i] = c.(network.Controller)
	}
	r.netRouter.LoadControllers(nc)

	for _, c := range controllers {
		baseSub := fmt.Sprintf(`%s.%s`, r.natsClient.Service.Info().Name, strings.ReplaceAll(c.Path(), "/", ""))

		ng := r.natsClient.Service.AddGroup(baseSub)
		c.MountNats(ng)
	}
}

func (r *router) Start(ip string, port uint16) {
	r.netRouter.Start(ip, port)
}

func (r *router) RegisterValidationParsers(tagNameFunc validator.TagNameFunc) {
	r.netRouter.RegisterValidationParsers(tagNameFunc)
}
