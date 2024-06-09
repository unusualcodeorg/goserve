package schema

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CollectionName = "messages"

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" validate:"-"`
	Type      string             `bson:"type" validate:"required"`
	Msg       string             `bson:"msg" validate:"required"`
	Status    bool               `bson:"status" validate:"required"`
	CreatedAt time.Time          `bson:"createdAt" validate:"required"`
	UpdatedAt time.Time          `bson:"updatedAt" validate:"required"`
}

func NewMessage(msgType string, msgTxt string) (mongo.Schema[Message], error) {
	time := time.Now()
	m := Message{
		Type:      msgType,
		Msg:       msgTxt,
		Status:    true,
		CreatedAt: time,
		UpdatedAt: time,
	}
	if err := m.Validate(); err != nil {
		return nil, err
	}
	return &m, nil
}

func (message *Message) GetDocument() *Message {
	return message
}

func (message *Message) Validate() error {
	validate := validator.New()
	return validate.Struct(message)
}

func (*Message) EnsureIndexes(db mongo.Database) {

}
