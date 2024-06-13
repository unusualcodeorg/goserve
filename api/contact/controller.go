package contact

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/contact/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/user/model"
	coredto "github.com/unusualcodeorg/go-lang-backend-architecture/core/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
	"github.com/unusualcodeorg/go-lang-backend-architecture/utils"
)

type controller struct {
	network.BaseController
	service Service
}

func NewContactController(
	authProvider network.AuthenticationProvider,
	authorizeProvider network.AuthorizationProvider,
	service Service,
) network.Controller {
	c := controller{
		BaseController: network.NewBaseController("/contact", authProvider, authorizeProvider),
		service: service,
	}
	return &c
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.Use(c.Authentication())
	group.Use(c.Authorization(string(model.RoleCodeLearner)))
	group.POST("/", c.createMessageHandler)
	group.GET("/id/:id", c.getMessageHandler)
	group.GET("/paginated", c.getMessagesPaginated)
}

func (c *controller) createMessageHandler(ctx *gin.Context) {
	body, err := network.ReqBody(ctx, &dto.CreateMessage{})
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	msg, err := c.service.SaveMessage(body)
	if err != nil {
		c.Send(ctx).InternalServerError("something went wrong", err)
		return
	}

	data, err := utils.MapTo[dto.InfoMessage](msg)
	if err != nil {
		c.Send(ctx).InternalServerError("something went wrong", err)
		return
	}

	c.Send(ctx).SuccessDataResponse("message received successfully!", data)
}

func (c *controller) getMessageHandler(ctx *gin.Context) {
	mongoId, err := network.ReqParams(ctx, &coredto.MongoId{})
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	msg, err := c.service.FindMessage(mongoId.ID)
	if err != nil {
		c.Send(ctx).NotFoundError("message not found", err)
		return
	}

	data, err := utils.MapTo[dto.InfoMessage](msg)
	if err != nil {
		c.Send(ctx).InternalServerError("something went wrong", err)
		return
	}

	c.Send(ctx).SuccessDataResponse("success", data)
}

func (c *controller) getMessagesPaginated(ctx *gin.Context) {
	pagination, err := network.ReqQuery(ctx, &coredto.Pagination{})
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	msgs, err := c.service.FindPaginatedMessage(pagination)

	if err != nil {
		c.Send(ctx).NotFoundError("messages not found", err)
		return
	}

	data, err := utils.MapTo[[]dto.InfoMessage](&msgs)
	if err != nil {
		c.Send(ctx).InternalServerError("something went wrong", err)
		return
	}

	c.Send(ctx).SuccessDataResponse("success", data)
}
