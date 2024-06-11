package user

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/profile/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

type controller struct {
	network.BaseController
	userService UserService
}

func NewUserController(
	authProvider network.AuthenticationProvider,
	authorizeProvider network.AuthorizationProvider,
	service UserService,
) network.Controller {
	c := controller{
		BaseController: network.NewBaseController("/profile", authProvider, authorizeProvider),
		userService:    service,
	}
	return &c
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.GET("/id/:id", c.getUserHandler)
}

func (c *controller) getUserHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	objectId, err := mongo.NewObjectID(id)
	if err != nil {
		panic(network.BadRequestError(err.Error(), err))
	}

	msg, err := c.userService.FindUserById(objectId)
	if err != nil {
		panic(network.NotFoundError("message not found", err))
	}

	data, err := network.MapToDto[dto.InfoPrivateUser](msg)
	if err != nil {
		panic(network.InternalServerError("something went wrong", err))
	}

	network.ResSuccessData(ctx, "success", data)
}
