package schema

import (
	"time"
)

const MessageCollectionName = "messages"

type Message struct {
	ID        string    `bson:"_id,omitempty" validate:"-"`
	Type      string    `bson:"type" validate:"required"`
	Msg       string    `bson:"msg" validate:"required"`
	Status    bool      `bson:"status" validate:"required"`
	CreatedAt time.Time `bson:"createdAt" validate:"required"`
	UpdatedAt time.Time `bson:"updatedAt" validate:"required"`
}

func NewMessage(msgType string, msgTxt string) *Message {
	time := time.Now()
	m := Message{
		Type:      msgType,
		Msg:       msgTxt,
		Status:    true,
		CreatedAt: time,
		UpdatedAt: time,
	}
	return &m
}
