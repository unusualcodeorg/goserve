package blog

import (
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/blog/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/blog/model"
	userModel "github.com/unusualcodeorg/go-lang-backend-architecture/api/user/model"
	"github.com/unusualcodeorg/go-lang-backend-architecture/framework/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/framework/network"
	"github.com/unusualcodeorg/go-lang-backend-architecture/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	GetPublisedBlogById(id primitive.ObjectID) (*dto.PublicBlog, error)
	GetPublishedBlogBySlug(slug string) (*dto.PublicBlog, error)
	CreateBlog(createBlogDto *dto.CreateBlog, author *userModel.User) (*dto.PrivateBlog, error)
}

type service struct {
	network.BaseService
	blogQueryBuilder mongo.QueryBuilder[model.Blog]
}

func NewService(db mongo.Database) Service {
	s := service{
		BaseService:      network.NewBaseService(),
		blogQueryBuilder: mongo.NewQueryBuilder[model.Blog](db, model.CollectionName),
	}
	return &s
}

func (s *service) GetPublisedBlogById(id primitive.ObjectID) (*dto.PublicBlog, error) {
	filter := bson.M{"_id": id, "isPublished": true, "status": true}

	blog, err := s.blogQueryBuilder.SingleQuery().FindOne(filter, nil)
	if err != nil {
		return nil, err
	}

	return utils.MapTo[dto.PublicBlog](blog)
}

func (s *service) GetPublishedBlogBySlug(slug string) (*dto.PublicBlog, error) {
	filter := bson.M{"slug": slug, "isPublished": true, "status": true}

	blog, err := s.blogQueryBuilder.SingleQuery().FindOne(filter, nil)
	if err != nil {
		return nil, err
	}

	return utils.MapTo[dto.PublicBlog](blog)
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

	return dto.NewPrivateBlog(created, author)
}
