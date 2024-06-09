package user

import (
	"time"

	"github.com/unusualcodeorg/go-lang-backend-architecture/api/user/schema"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService interface {
	FindUser(id *primitive.ObjectID) (*schema.User, error)
}

type service struct {
	network.BaseService
	userQuery mongo.DatabaseQuery[schema.User]
}

func NewUserService(db mongo.Database, dbQueryTimeout time.Duration) UserService {
	s := service{
		BaseService:  network.NewBaseService(dbQueryTimeout),
		userQuery: mongo.NewDatabaseQuery[schema.User](db, schema.CollectionName),
	}
	return &s
}

func (s *service) FindUser(id *primitive.ObjectID) (*schema.User, error) {
	ctx, cancel := s.Context()
	defer cancel()

	filter := bson.M{"_id": id}

	msg, err := s.userQuery.FindOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
