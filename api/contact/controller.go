package contact

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/contact/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/user/model"
	coredto "github.com/unusualcodeorg/go-lang-backend-architecture/core/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

type controller struct {
	network.BaseController
	contactService ContactService
}

func NewContactController(
	authProvider network.AuthenticationProvider,
	authorizeProvider network.AuthorizationProvider,
	service ContactService,
) network.Controller {
	c := controller{
		BaseController: network.NewBaseController("/contact", authProvider, authorizeProvider),
		contactService: service,
	}
	return &c
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.Use(c.Authentication())
	group.Use(c.Authorization(string(model.RoleCodeWriter)))
	group.POST("/", c.createMessageHandler)
	group.GET("/id/:id", c.getMessageHandler)
	group.GET("/paginated", c.getMessagesPaginated)
}

func (c *controller) createMessageHandler(ctx *gin.Context) {
	body, err := network.ReqBody(ctx, &dto.CreateMessage{})
	if err != nil {
		panic(network.BadRequestError(err.Error(), err))
	}

	msg, err := c.contactService.SaveMessage(body)
	if err != nil {
		panic(network.InternalServerError("something went wrong", err))
	}

	data, err := network.MapToDto[dto.InfoMessage](msg)
	if err != nil {
		panic(network.InternalServerError("something went wrong", err))
	}

	network.SuccessDataResponse(ctx, "message received successfully!", data)
}

func (c *controller) getMessageHandler(ctx *gin.Context) {
	mongoId, err := network.ReqParams(ctx, &coredto.MongoId{})
	if err != nil {
		panic(network.BadRequestError(err.Error(), err))
	}

	msg, err := c.contactService.FindMessage(mongoId.ID)
	if err != nil {
		panic(network.NotFoundError("message not found", err))
	}

	data, err := network.MapToDto[dto.InfoMessage](msg)
	if err != nil {
		panic(network.InternalServerError("something went wrong", err))
	}

	network.SuccessDataResponse(ctx, "success", data)
}

func (c *controller) getMessagesPaginated(ctx *gin.Context) {
	pagination, err := network.ReqQuery(ctx, &coredto.Pagination{})
	if err != nil {
		panic(network.BadRequestError(err.Error(), err))
	}

	msgs, err := c.contactService.FindPaginatedMessage(pagination)

	if err != nil {
		panic(network.NotFoundError("messages not found", err))
	}

	data, err := network.MapToDto[[]dto.InfoMessage](&msgs)
	if err != nil {
		panic(network.InternalServerError("something went wrong", err))
	}
	network.SuccessDataResponse(ctx, "success", data)
}
