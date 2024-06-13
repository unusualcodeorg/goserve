package blog

import (
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/blog/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/blog/model"
	userModel "github.com/unusualcodeorg/go-lang-backend-architecture/api/user/model"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogService interface {
	// GetBlog(id primitive.ObjectID) (*dto.InfoBlog, error)
	CreateBlog(createBlogDto *dto.CreateBlog, author *userModel.User) (*dto.PrivateBlog, error)
}

type service struct {
	network.BaseService
	blogQueryBuilder mongo.QueryBuilder[model.Blog]
}

func NewBlogService(db mongo.Database) BlogService {
	s := service{
		BaseService:      network.NewBaseService(),
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

func (s *service) CreateBlog(b *dto.CreateBlog, author *userModel.User) (*dto.PrivateBlog, error) {
	filter := bson.M{"slug": b.Slug}
	_, err := s.blogQueryBuilder.SingleQuery().FindOne(filter, nil)
	if err == nil {
		return nil, network.NewBadRequestError("Blog with slug: "+b.Slug+" already exists", nil)
	}

	blog, err := model.NewBlog(b.Slug, b.Title, b.Description, b.DraftText, b.Tags, author)
	if err != nil {
		return nil, err
	}

	created, err := s.blogQueryBuilder.SingleQuery().InsertAndRetrieveOne(blog)
	if err != nil {
		return nil, err
	}

	return dto.NewPrivateBlog(created)
}
