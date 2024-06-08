package contact

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/contact/dto"
	coredto "github.com/unusualcodeorg/go-lang-backend-architecture/core/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

type controller struct {
	network.BaseController
	contactService ContactService
}

func NewContactController(
	authMFunc network.GroupMiddlewareFunc,
	authorizeMFunc network.GroupMiddlewareFunc,
	service ContactService,
) network.Controller {
	c := controller{
		BaseController: network.NewBaseController("/contact", authMFunc, authorizeMFunc),
		contactService: service,
	}
	return &c
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.POST("/", c.createMessageHandler)
	group.GET("/id/:id", c.getMessageHandler)
	group.GET("/paginated", c.getMessagesPaginated)
}

func (c *controller) createMessageHandler(ctx *gin.Context) {
	body, err := network.ReqBody[dto.CreateMessage](ctx)
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

	network.SuccessResponse("message received successfully!", data).Send(ctx)
}

func (c *controller) getMessageHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	objectId, err := mongo.NewObjectID(id)
	if err != nil {
		panic(network.BadRequestError(err.Error(), err))
	}

	msg, err := c.contactService.FindMessage(objectId)
	if err != nil {
		panic(network.NotFoundError("message not found", err))
	}

	data, err := network.MapToDto[dto.InfoMessage](msg)
	if err != nil {
		panic(network.InternalServerError("something went wrong", err))
	}

	network.SuccessResponse("success", data).Send(ctx)
}

func (c *controller) getMessagesPaginated(ctx *gin.Context) {
	pagenation, err := network.ReqQuery[coredto.PaginationDto](ctx)
	if err != nil {
		panic(network.BadRequestError(err.Error(), err))
	}

	msgs, err := c.contactService.FindPaginatedMessage(pagenation)

	if err != nil {
		panic(network.NotFoundError("messages not found", err))
	}

	data, err := network.MapToDto[[]dto.InfoMessage](msgs)
	if err != nil {
		panic(network.InternalServerError("something went wrong", err))
	}
	network.SuccessResponse("success", data).Send(ctx)
}
