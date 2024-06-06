package api

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/contact"
	"github.com/unusualcodeorg/go-lang-backend-architecture/internal/handlers"
)

func CreateRouter() *gin.Engine {
	router := gin.Default()
	loadControllers(router)
	router.NoRoute(handlers.NotFoundHandler)
	return router
}

func loadControllers(router *gin.Engine) {
	contact.Controller(router)
}
