package contact

import (
	"context"
	"time"

	"github.com/unusualcodeorg/go-lang-backend-architecture/api/contact/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/contact/schema"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
)

type ContactService interface {
	SaveMessage(d *dto.CreateMessage) (*schema.Message, error)
	FindMessage(id *mongo.ObjectID) (*schema.Message, error)
	FindPaginatedMessage(page uint64, limit uint64) (*[]schema.Message, error)
}

type service struct {
	messageDbQuery mongo.DatabaseQuery[schema.Message]
}

func NewContactService(db mongo.Database) ContactService {
	s := service{
		messageDbQuery: mongo.NewDatabaseQuery[schema.Message](db, schema.MessageCollectionName),
	}
	return &s
}

func (s *service) SaveMessage(d *dto.CreateMessage) (*schema.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	msg, err := schema.NewMessage(d.Type, d.Msg)
	if err != nil {
		return nil, err
	}

	result, err := s.messageDbQuery.InsertOne(ctx, msg)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) FindMessage(id *mongo.ObjectID) (*schema.Message, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := mongo.Filter{"_id": id.ObjectID}

	msg, err := s.messageDbQuery.FindOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (s *service) FindPaginatedMessage(page uint64, limit uint64) (*[]schema.Message, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := mongo.Filter{"status": true}

	msgs, err := s.messageDbQuery.FindPaginated(ctx, filter, 2, 10)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}
