package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/auth/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

type controller struct {
	network.BaseController
	authService AuthService
}

func NewAuthController(
	authProvider network.AuthenticationProvider,
	authorizeProvider network.AuthorizationProvider,
	service AuthService,
) network.Controller {
	c := controller{
		BaseController: network.NewBaseController("/auth", authProvider, authorizeProvider),
		authService:    service,
	}
	return &c
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.POST("/signup/basic", c.signUpBasicHandler)
	group.POST("/signin/basic", c.signInBasicHandler)
	group.DELETE("/signout", c.Authentication(), c.signOutBasicHandler)
}

func (c *controller) signUpBasicHandler(ctx *gin.Context) {
	body, err := network.ReqBody(ctx, dto.EmptySignUpBasic())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	exists := c.authService.IsEmailRegisted(body.Email)
	if exists {
		c.Send(ctx).BadRequestError("user already exists", nil)
		return
	}

	data, err := c.authService.SignUpBasic(body)
	if err != nil {
		c.Send(ctx).InternalServerError(err.Error(), err)
		return
	}

	c.Send(ctx).SuccessDataResponse("success", data)
}

func (c *controller) signInBasicHandler(ctx *gin.Context) {
	body, err := network.ReqBody(ctx, dto.EmptySignInBasic())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	exists := c.authService.IsEmailRegisted(body.Email)
	if !exists {
		c.Send(ctx).NotFoundError("user not registered", nil)
		return
	}

	// bcrypt.CompareHashAndPassword()
}

func (c *controller) signOutBasicHandler(ctx *gin.Context) {
	c.Send(ctx).SuccessMsgResponse("logout not working!")
}
