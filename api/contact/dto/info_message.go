package dto

import (
	"time"

	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
)

type InfoMessage struct {
	ID        mongo.ObjectID `json:"_id" binding:"required"`
	Type      string         `json:"type" binding:"required"`
	Msg       string         `json:"msg" binding:"required"`
	CreatedAt time.Time      `json:"createdAt" binding:"required"`
}
