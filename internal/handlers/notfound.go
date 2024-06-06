package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/internal/core"
)

func NotFoundHandler(ctx *gin.Context) {
	core.NotFoundResponse("resource not found").Send(ctx)
}
