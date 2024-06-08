package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InfoMessage struct {
	ID        primitive.ObjectID `json:"_id" binding:"required"`
	Type      string             `json:"type" binding:"required"`
	Msg       string             `json:"msg" binding:"required"`
	CreatedAt time.Time          `json:"createdAt" binding:"required"`
}
