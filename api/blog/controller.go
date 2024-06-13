package blog

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

type controller struct {
	network.BaseController
	blogService BlogService
}

func NewController(
	authMFunc network.AuthenticationProvider,
	authorizeMFunc network.AuthorizationProvider,
	service BlogService,
) network.Controller {
	return &controller{
		BaseController: network.NewBaseController("/blog", authMFunc, authorizeMFunc),
		blogService:    service,
	}
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.GET("/id/:id", c.getBlogHandler)

}

func (c *controller) getBlogHandler(ctx *gin.Context) {

}
