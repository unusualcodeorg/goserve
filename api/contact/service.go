package contact

import (
	"context"

	"github.com/unusualcodeorg/go-lang-backend-architecture/api/contact/schema"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ContactService interface {
	SaveMessage(msgType string, msgTxt string) (*schema.Message, error)
}

type service struct {
	db mongo.Database
}

func NewService(database mongo.Database) ContactService {
	s := service{
		db: database,
	}
	return &s
}

func (s *service) SaveMessage(msgType string, msgTxt string) (*schema.Message, error) {

	msg := schema.NewMessage(msgType, msgTxt)

	if err := utils.Validate(msg); err != nil {
		return nil, err
	}

	collection := s.db.GetCollection(schema.MessageCollectionName)

	result, err := collection.InsertOne(context.Background(), msg)
	if err != nil {
		return nil, err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, err
	}

	msg.ID = insertedID.Hex()
	return msg, nil
}
