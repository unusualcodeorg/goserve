package message

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/message/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/internal/core"
)

func Controller(router *gin.Engine) {
	router.POST("/message", createMessageHandler)
}

func createMessageHandler(ctx *gin.Context) {
	var message dto.CreateMessage

	if err := core.ParseBody(ctx, &message); err != nil {
		core.BadRequestResponse(err).Send(ctx)
		return
	}

	core.SuccessResponse(message).Send(ctx)
}
