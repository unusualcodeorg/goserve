package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id" binding:"required,uuid"`    // ID must be a valid UUID
	Title  string  `json:"title" binding:"required"`      // Title is required
	Artist string  `json:"artist" binding:"required"`     // Artist is required
	Price  float64 `json:"price" binding:"required,gt=0"` // Price is required and must be greater than 0
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	router.Run("localhost:8080")
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.JSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum album

	// Bind JSON and validate
	if err := c.ShouldBindJSON(&newAlbum); err != nil {
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
			c.JSON(http.StatusBadRequest, gin.H{"errors": errorMessages})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.JSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range albums {
		if a.ID == id {
			c.JSON(http.StatusOK, a)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
