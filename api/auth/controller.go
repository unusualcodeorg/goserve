package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/auth/dto"
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
	group.POST("/signup/basic", c.singupBasicHandler)

	logout := group.Group("/logout")
	c.AuthenticationMiddleware().Attach(logout)
	logout.DELETE("/", c.logoutBasicHandler)

}

func (c *controller) singupBasicHandler(ctx *gin.Context) {
	body, err := network.ReqBody(ctx, dto.EmptySignUpBasic())
	if err != nil {
		panic(network.BadRequestError(err.Error(), err))
	}

	exists := c.authService.IsEmailRegisted(body.Email)
	if exists {
		e := errors.New("user already exists")
		panic(network.BadRequestError(e.Error(), e))
	}

	data, err := c.authService.SignUpBasic(body)

	if err != nil {
		panic(network.InternalServerError(err.Error(), err))
	}

	network.SuccessResponse("success", data).Send(ctx)
}

func (c *controller) logoutBasicHandler(ctx *gin.Context) {
	network.SuccessMsgResponse("logout not working!").Send(ctx)
}
