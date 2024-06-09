package auth

import (
	"time"

	"github.com/unusualcodeorg/go-lang-backend-architecture/api/auth/schema"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthService interface {
	FindAuth(id *primitive.ObjectID) (*schema.Role, error)
}

type service struct {
	network.BaseService
	authQuery mongo.DatabaseQuery[schema.Role]
}

func NewAuthService(db mongo.Database, dbQueryTimeout time.Duration) AuthService {
	s := service{
		BaseService:  network.NewBaseService(dbQueryTimeout),
		authQuery: mongo.NewDatabaseQuery[schema.Role](db, schema.CollectionName),
	}
	return &s
}

func (s *service) FindAuth(id *primitive.ObjectID) (*schema.Role, error) {
	ctx, cancel := s.Context()
	defer cancel()

	filter := bson.M{"_id": id}

	msg, err := s.authQuery.FindOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
