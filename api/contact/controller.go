package contact

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/contact/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/utils"
)

type controller struct {
	network.BaseController
	contactService ContactService
}

func NewContactController(
	base network.BaseController,
	s ContactService,
) network.Controller {
	c := controller{
		BaseController: base,
		contactService: s,
	}
	return &c
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.POST("/", c.createMessageHandler)
}

func (c *controller) createMessageHandler(ctx *gin.Context) {
	var createMsg dto.CreateMessage

	if err := utils.GetBody(ctx, &createMsg); err != nil {
		network.BadRequestResponse(err).Send(ctx)
		return
	}

	_, err := c.contactService.SaveMessage(createMsg.Type, createMsg.Msg)

	if err != nil {
		network.InternalServerErrorResponse("Something went wrong")
		return
	}

	network.SuccessMsgResponse("Message received successfully!").Send(ctx)
}
