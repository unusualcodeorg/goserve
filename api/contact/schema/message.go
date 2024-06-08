package schema

import (
	"time"

	"github.com/unusualcodeorg/go-lang-backend-architecture/core/parser"
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

func NewMessage(msgType string, msgTxt string) (*Message, error) {
	time := time.Now()
	m := Message{
		Type:      msgType,
		Msg:       msgTxt,
		Status:    true,
		CreatedAt: time,
		UpdatedAt: time,
	}
	if err := parser.Validate(m); err != nil {
		return nil, err
	}
	return &m, nil
}
