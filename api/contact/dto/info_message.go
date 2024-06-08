package dto

import "time"

type InfoMessage struct {
	Type      string    `json:"type" binding:"required"`
	Msg       string    `json:"msg" binding:"required"`
	CreatedAt time.Time `json:"createdAt" binding:"required"`
}
