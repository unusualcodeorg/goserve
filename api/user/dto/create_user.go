package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InfoUser struct {
	ID        primitive.ObjectID `json:"_id" binding:"required"`
	Field     string             `json:"field" binding:"required"`
	CreatedAt time.Time          `json:"createdAt" binding:"required"`
}
