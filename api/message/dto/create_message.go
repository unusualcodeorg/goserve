package dto

type CreateMessage struct {
	Type    string `json:"type" binding:"required"`
	Message string `json:"message" binding:"required"`
}
