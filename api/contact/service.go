package contact

import (
	"context"
	"time"

	"github.com/unusualcodeorg/go-lang-backend-architecture/api/contact/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/contact/schema"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ContactService interface {
	SaveMessage(d dto.CreateMessage) (*schema.Message, error)
	FindMessage(id primitive.ObjectID) (*schema.Message, error)
}

type service struct {
	dbQuery mongo.DatabaseQuery
}

func NewContactService(dbQuery mongo.DatabaseQuery) ContactService {
	s := service{
		dbQuery: dbQuery,
	}
	return &s
}

func (s *service) SaveMessage(d dto.CreateMessage) (*schema.Message, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	msg, err := schema.NewMessage(d.Type, d.Msg)
	if err != nil {
		return nil, err
	}

	result, err := s.dbQuery.InsertOne(ctx, schema.MessageCollectionName, msg)
	if err != nil {
		return nil, err
	}

	msg.ID = result.Hex()

	return msg, nil
}

func (s *service) FindMessage(id primitive.ObjectID) (*schema.Message, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var msg schema.Message
	filter := bson.M{"_id": id}

	err := s.dbQuery.FindOne(ctx, schema.MessageCollectionName, filter, &msg)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}
