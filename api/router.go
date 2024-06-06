package api

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/contact"
)

func CreateRouter() *gin.Engine {
	router := gin.Default()
	loadControllers(router)
	return router
}

func loadControllers(router *gin.Engine) {
	contact.Controller(router)
}
