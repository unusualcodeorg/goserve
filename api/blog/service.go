package blog

import (
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
	blogQueryBuilder mongo.QueryBuilder[model.Blog]
}

func NewBlogService(db mongo.Database) BlogService {
	s := service{
		BaseService:  network.NewBaseService(),
		blogQueryBuilder: mongo.NewQueryBuilder[model.Blog](db, model.CollectionName),
	}
	return &s
}

func (s *service) FindBlog(id primitive.ObjectID) (*model.Blog, error) {
	filter := bson.M{"_id": id}

	msg, err := s.blogQueryBuilder.SingleQuery().FindOne(filter, nil)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
