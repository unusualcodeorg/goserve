package contact

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/contact/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/middleware"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type controller struct {
	network.BaseController
	contactService ContactService
}

func NewContactController(service ContactService) network.Controller {
	path := "/contact"
	base := network.NewBaseController(path, middleware.NewAuthentication, middleware.NewAuthorization)
	c := controller{
		BaseController: base,
		contactService: service,
	}
	return &c
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.POST("/", c.createMessageHandler)
	group.GET("/id/:id", c.getMessageHandler)
}

func (c *controller) createMessageHandler(ctx *gin.Context) {
	var body dto.CreateMessage

	if err := network.Body(ctx, &body); err != nil {
		network.BadRequestResponse(err).Send(ctx)
		return
	}

	_, err := c.contactService.SaveMessage(body)

	if err != nil {
		network.InternalServerErrorResponse("something went wrong")
		return
	}

	network.SuccessMsgResponse("message received successfully!").Send(ctx)
}

func (c *controller) getMessageHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		network.BadRequestResponse([]string{id + " is not a valid mongo id"}).Send(ctx)
		return
	}

	msg, err := c.contactService.FindMessage(objectId)

	if err != nil {
		network.NotFoundResponse("message not found").Send(ctx)
		return
	}

	data := network.MapToDto(msg, &dto.InfoMessage{})
	network.SuccessResponse("success", data).Send(ctx)
}
