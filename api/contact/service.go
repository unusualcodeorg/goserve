package contact

import (
	"time"

	"github.com/unusualcodeorg/go-lang-backend-architecture/api/contact/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/contact/schema"
	coredto "github.com/unusualcodeorg/go-lang-backend-architecture/core/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ContactService interface {
	SaveMessage(d *dto.CreateMessage) (*schema.Message, error)
	FindMessage(id *primitive.ObjectID) (*schema.Message, error)
	FindPaginatedMessage(p *coredto.PaginationDto) ([]schema.Message, error)
}

type service struct {
	network.BaseService
	messageQuery mongo.DatabaseQuery[schema.Message]
}

func NewContactService(db mongo.Database, dbQueryTimeout time.Duration) ContactService {
	s := service{
		BaseService:  network.NewBaseService(dbQueryTimeout),
		messageQuery: mongo.NewDatabaseQuery[schema.Message](db, schema.CollectionName),
	}
	return &s
}

func (s *service) SaveMessage(d *dto.CreateMessage) (*schema.Message, error) {
	ctx, cancel := s.Context()
	defer cancel()

	msg, err := schema.NewMessage(d.Type, d.Msg)
	if err != nil {
		return nil, err
	}

	result, err := s.messageQuery.InsertAndRetrieveOne(ctx, msg)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) FindMessage(id *primitive.ObjectID) (*schema.Message, error) {
	ctx, cancel := s.Context()
	defer cancel()

	filter := bson.M{"_id": id}

	msg, err := s.messageQuery.FindOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (s *service) FindPaginatedMessage(p *coredto.PaginationDto) ([]schema.Message, error) {
	ctx, cancel := s.Context()
	defer cancel()

	filter := bson.M{"status": true}

	msgs, err := s.messageQuery.FindPaginated(ctx, filter, p.Page, p.Limit)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}
