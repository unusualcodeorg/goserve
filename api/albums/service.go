package albums

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/unusualcodeorg/go-lang-backend-architecture/api/albums/dto"
)

// getAlbums responds with the list of all albums as JSON.
func GetAlbumsHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, AlbumStore)
}

// postAlbums adds an album from JSON received in the request body.
func PostAlbumsHandler(ctx *gin.Context) {
	var newAlbum dto.Album

	// Bind JSON and validate
	if err := ctx.ShouldBindJSON(&newAlbum); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMessages := make([]string, 0)
			for _, err := range validationErrors {
				switch err.Tag() {
				case "required":
					errorMessages = append(errorMessages, err.Field()+" is required")
				case "uuid":
					errorMessages = append(errorMessages, err.Field()+" must be a valid UUID")
				case "gt":
					errorMessages = append(errorMessages, err.Field()+" must be greater than "+err.Param())
				default:
					errorMessages = append(errorMessages, err.Field()+" is invalid")
				}
			}
			ctx.JSON(http.StatusBadRequest, gin.H{"errors": errorMessages})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Add the new album to the slice.
	AlbumStore = append(AlbumStore, newAlbum)
	ctx.JSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func GetAlbumByIDHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range AlbumStore {
		if a.ID == id {
			ctx.JSON(http.StatusOK, a)
			return
		}
	}
	ctx.JSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
