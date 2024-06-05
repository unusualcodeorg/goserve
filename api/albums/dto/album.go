package dto

type Album struct {
	ID     string  `json:"id" binding:"required,gt=0"`    // ID is required
	Title  string  `json:"title" binding:"required"`      // Title is required
	Artist string  `json:"artist" binding:"required"`     // Artist is required
	Price  float64 `json:"price" binding:"required,gt=0"` // Price is required and must be greater than 0
}
