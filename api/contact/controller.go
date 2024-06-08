package contact

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/contact/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
	"github.com/unusualcodeorg/go-lang-backend-architecture/utils"
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
	var body dto.CreateMessage

	if err := utils.Body(ctx, &body); err != nil {
		network.BadRequestResponse(err).Send(ctx)
		return
	}

	msg, err := c.contactService.SaveMessage(&body)
	if err != nil {
		network.InternalServerErrorResponse("something went wrong")
		return
	}

	data, err := utils.MapToDto[dto.InfoMessage](msg)
	if err != nil {
		network.InternalServerErrorResponse("something went wrong")
		return
	}

	network.SuccessResponse("message received successfully!", data).Send(ctx)
}

func (c *controller) getMessageHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	objectId, err := mongo.NewObjectID(id)
	if err != nil {
		network.BadRequestResponse([]string{id + " is not a valid mongo id"}).Send(ctx)
		return
	}

	msg, err := c.contactService.FindMessage(objectId)

	if err != nil {
		network.NotFoundResponse("message not found").Send(ctx)
		return
	}

	data, err := utils.MapToDto[dto.InfoMessage](msg)
	if err != nil {
		network.InternalServerErrorResponse("something went wrong")
		return
	}
	network.SuccessResponse("success", data).Send(ctx)
}

func (c *controller) getMessagesPaginated(ctx *gin.Context) {

	msgs, err := c.contactService.FindPaginatedMessage(1, 5)

	if err != nil {
		network.NotFoundResponse("message not found").Send(ctx)
		return
	}

	data, err := utils.MapToDto[[]dto.InfoMessage](msgs)
	if err != nil {
		network.InternalServerErrorResponse("something went wrong")
		return
	}
	network.SuccessResponse("success", data).Send(ctx)
}
