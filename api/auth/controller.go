package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

type controller struct {
	network.BaseController
	authService AuthService
}

func NewAuthController(
	authMFunc network.GroupMiddlewareFunc,
	authorizeMFunc network.GroupMiddlewareFunc,
	service AuthService,
) network.Controller {
	c := controller{
		BaseController: network.NewBaseController("/auth", authMFunc, authorizeMFunc),
		authService:    service,
	}
	return &c
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.GET("/id/:id", c.getAuthHandler)
}

func (c *controller) getAuthHandler(ctx *gin.Context) {

	network.SuccessResponse("success", "").Send(ctx)
}
