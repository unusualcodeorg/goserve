package schema

import (
	"time"
)

const collection_name = "messages"

type Message struct {
	CollectionName string    `bson:"-" validate:"-"`
	ID             string    `bson:"_id,omitempty" validate:"-"`
	Type           string    `bson:"type" validate:"required"`
	Msg            string    `bson:"msg" validate:"required"`
	Status         bool      `bson:"status" validate:"required"`
	CreatedAt      time.Time `bson:"createdAt" validate:"required"`
	UpdatedAt      time.Time `bson:"updatedAt" validate:"required"`
}

func NewMessage(msgType string, msgTxt string) *Message {
	time := time.Now()
	return &Message{
		CollectionName: collection_name,
		Type:           msgType,
		Msg:            msgTxt,
		Status:         true,
		CreatedAt:      time,
		UpdatedAt:      time,
	}
}
