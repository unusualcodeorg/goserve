package contact

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/contact/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/internal/core"
)

func Controller(router *gin.Engine) {
	router.POST("/message", createMessageHandler)
}

func createMessageHandler(ctx *gin.Context) {
	var createMsg dto.CreateMessage

	if err := core.ParseBody(ctx, &createMsg); err != nil {
		core.BadRequestResponse(err).Send(ctx)
		return
	}

	_, err := saveMessage(createMsg.Type, createMsg.Msg)

	if err != nil {
		core.InternalServerErrorResponse("Something went wrong")
		return
	}

	core.SuccessMsgResponse("Message received successfully!").Send(ctx)
}
