package network

import (
	"github.com/gin-gonic/gin"
)

func NotFoundHandler(ctx *gin.Context) {
	NotFoundResponse("resource not found").Send(ctx)
}
