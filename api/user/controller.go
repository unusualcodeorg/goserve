package user

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/user/dto"
	coredto "github.com/unusualcodeorg/go-lang-backend-architecture/core/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

type controller struct {
	network.BaseController
	userService UserService
}

func NewProfileController(
	authProvider network.AuthenticationProvider,
	authorizeProvider network.AuthorizationProvider,
	userService UserService,
) network.Controller {
	return &controller{
		BaseController: network.NewBaseController("/profile", authProvider, authorizeProvider),
		userService:    userService,
	}
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.GET("/id/:id", c.getUserHandler)
}

func (c *controller) getUserHandler(ctx *gin.Context) {
	mongoId, err := network.ReqParams(ctx, &coredto.MongoId{})
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	msg, err := c.userService.FindUserById(mongoId.ID)
	if err != nil {
		c.Send(ctx).NotFoundError("message not found", err)
		return
	}

	data, err := network.MapToDto[dto.InfoPrivateUser](msg)
	if err != nil {
		c.Send(ctx).InternalServerError("something went wrong", err)
		return
	}

	c.Send(ctx).SuccessDataResponse("success", data)
}
