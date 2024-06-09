package core

import (
	"time"

	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/schema"
	"go.mongodb.org/mongo-driver/bson"
)

type CoreService interface {
	FindApiKey(key string) (*schema.ApiKey, error)
}

type service struct {
	network.BaseService
	apikeyQuery mongo.Query[schema.ApiKey]
}

func NewCoreService(db mongo.Database, dbQueryTimeout time.Duration) CoreService {
	s := service{
		BaseService: network.NewBaseService(dbQueryTimeout),
		apikeyQuery: mongo.NewQuery[schema.ApiKey](db, schema.CollectionName),
	}
	return &s
}

func (s *service) FindApiKey(key string) (*schema.ApiKey, error) {
	ctx, cancel := s.Context()
	defer cancel()

	filter := bson.M{"key": key, "status": true}

	apikey, err := s.apikeyQuery.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, err
	}

	return apikey, nil
}
