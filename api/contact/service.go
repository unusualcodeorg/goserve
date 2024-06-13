package contact

import (
	"time"

	"github.com/unusualcodeorg/go-lang-backend-architecture/api/contact/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/contact/model"
	coredto "github.com/unusualcodeorg/go-lang-backend-architecture/core/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ContactService interface {
	SaveMessage(d *dto.CreateMessage) (*model.Message, error)
	FindMessage(id primitive.ObjectID) (*model.Message, error)
	FindPaginatedMessage(p *coredto.Pagination) ([]model.Message, error)
}

type service struct {
	network.BaseService
	messageQueryBuilder mongo.QueryBuilder[model.Message]
}

func NewContactService(db mongo.Database, dbQueryTimeout time.Duration) ContactService {
	s := service{
		BaseService:  network.NewBaseService(dbQueryTimeout),
		messageQueryBuilder: mongo.NewQueryBuilder[model.Message](db, model.CollectionName),
	}
	return &s
}

func (s *service) SaveMessage(d *dto.CreateMessage) (*model.Message, error) {
	ctx, cancel := s.Context()
	defer cancel()

	msg, err := model.NewMessage(d.Type, d.Msg)
	if err != nil {
		return nil, err
	}

	result, err := s.messageQueryBuilder.Query(ctx).InsertAndRetrieveOne(msg)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) FindMessage(id primitive.ObjectID) (*model.Message, error) {
	ctx, cancel := s.Context()
	defer cancel()

	filter := bson.M{"_id": id}

	msg, err := s.messageQueryBuilder.Query(ctx).FindOne(filter, nil)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (s *service) FindPaginatedMessage(p *coredto.Pagination) ([]model.Message, error) {
	ctx, cancel := s.Context()
	defer cancel()

	filter := bson.M{"status": true}

	msgs, err := s.messageQueryBuilder.Query(ctx).FindPaginated(filter, p.Page, p.Limit, nil)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}
