package blog

import (
	"time"

	"github.com/unusualcodeorg/go-lang-backend-architecture/api/blog/model"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogService interface {
	FindBlog(id primitive.ObjectID) (*model.Blog, error)
}

type service struct {
	network.BaseService
	blogQuery mongo.Query[model.Blog]
}

func NewBlogService(db mongo.Database, dbQueryTimeout time.Duration) BlogService {
	s := service{
		BaseService:  network.NewBaseService(dbQueryTimeout),
		blogQuery: mongo.NewQuery[model.Blog](db, model.CollectionName),
	}
	return &s
}

func (s *service) FindBlog(id primitive.ObjectID) (*model.Blog, error) {
	ctx, cancel := s.Context()
	defer cancel()

	filter := bson.M{"_id": id}

	msg, err := s.blogQuery.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
