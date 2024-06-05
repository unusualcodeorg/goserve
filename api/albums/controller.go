package albums

import (
	"github.com/gin-gonic/gin"
)

func Controller(router *gin.Engine) {
	router.GET("/albums", GetAlbumsHandler)
	router.GET("/albums/:id", GetAlbumByIDHandler)
	router.POST("/albums", PostAlbumsHandler)
}
