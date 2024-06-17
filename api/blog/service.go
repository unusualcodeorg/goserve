package blog

import (
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/blog/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/blog/model"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/user"
	userModel "github.com/unusualcodeorg/go-lang-backend-architecture/api/user/model"
	coredto "github.com/unusualcodeorg/go-lang-backend-architecture/framework/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/framework/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/framework/network"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	CreateBlog(createBlogDto *dto.CreateBlog, author *userModel.User) (*dto.PrivateBlog, error)
	BlogSubmission(blogId primitive.ObjectID, author *userModel.User, submit bool) error
	GetPrivateBlogById(id primitive.ObjectID, author *userModel.User) (*dto.PrivateBlog, error)
	GetPublisedBlogById(id primitive.ObjectID) (*dto.PublicBlog, error)
	GetPublishedBlogBySlug(slug string) (*dto.PublicBlog, error)
	GetPaginatedDraftsForAuthor(author *userModel.User, p *coredto.Pagination) ([]*dto.InfoBlog, error)
	GetPaginatedPublishedForAuthor(author *userModel.User, p *coredto.Pagination) ([]*dto.InfoBlog, error)
	GetPaginatedSubmittedForAuthor(author *userModel.User, p *coredto.Pagination) ([]*dto.InfoBlog, error)
}

type service struct {
	network.BaseService
	blogQueryBuilder mongo.QueryBuilder[model.Blog]
	userService      user.Service
}

func NewService(db mongo.Database) Service {
	s := service{
		BaseService:      network.NewBaseService(),
		blogQueryBuilder: mongo.NewQueryBuilder[model.Blog](db, model.CollectionName),
	}
	return &s
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

func (s *service) BlogSubmission(blogId primitive.ObjectID, author *userModel.User, submit bool) error {
	filter := bson.M{"_id": blogId, "author": author.ID, "status": true}
	update := bson.M{"$set": bson.M{"isSubmitted": submit}}
	result, err := s.blogQueryBuilder.SingleQuery().UpdateOne(filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		network.NewNotFoundError("blog not found", nil)
		return nil
	}

	return nil
}

func (s *service) GetPrivateBlogById(id primitive.ObjectID, author *userModel.User) (*dto.PrivateBlog, error) {
	filter := bson.M{"_id": id, "author": author.ID, "status": true}

	blog, err := s.blogQueryBuilder.SingleQuery().FindOne(filter, nil)
	if err != nil {
		return nil, err
	}

	return dto.NewPrivateBlog(blog, author)
}

func (s *service) GetPublisedBlogById(id primitive.ObjectID) (*dto.PublicBlog, error) {
	filter := bson.M{"_id": id, "isPublished": true, "status": true}
	return s.getPublicPublishedBlog(filter)
}

func (s *service) GetPublishedBlogBySlug(slug string) (*dto.PublicBlog, error) {
	filter := bson.M{"slug": slug, "isPublished": true, "status": true}
	return s.getPublicPublishedBlog(filter)
}

func (s *service) getPublicPublishedBlog(filter bson.M) (*dto.PublicBlog, error) {
	projection := bson.D{{Key: "text", Value: 0}, {Key: "draftText", Value: 0}, {Key: "text", Value: 0}}
	opts := options.FindOne().SetProjection(projection)
	blog, err := s.blogQueryBuilder.SingleQuery().FindOne(filter, opts)
	if err != nil {
		return nil, err
	}

	author, err := s.userService.FindUserPublicProfile(blog.Author)
	if err != nil {
		return nil, err
	}

	return dto.NewPublicBlog(blog, author)
}

func (s *service) GetPaginatedDraftsForAuthor(author *userModel.User, p *coredto.Pagination) ([]*dto.InfoBlog, error) {
	filter := bson.M{"author": author.ID, "status": true, "isDraft": true}
	return s.getPaginated(filter, p)
}

func (s *service) GetPaginatedPublishedForAuthor(author *userModel.User, p *coredto.Pagination) ([]*dto.InfoBlog, error) {
	filter := bson.M{"author": author.ID, "status": true, "isPublished": true}
	return s.getPaginated(filter, p)
}

func (s *service) GetPaginatedSubmittedForAuthor(author *userModel.User, p *coredto.Pagination) ([]*dto.InfoBlog, error) {
	filter := bson.M{"author": author.ID, "status": true, "isSubmitted": true}
	return s.getPaginated(filter, p)
}

func (s *service) getPaginated(filter bson.M, p *coredto.Pagination) ([]*dto.InfoBlog, error) {
	blogs, err := s.blogQueryBuilder.SingleQuery().FindPaginated(filter, p.Page, p.Limit, nil)
	if err != nil {
		return nil, err
	}

	dtos := make([]*dto.InfoBlog, len(blogs))

	for i, b := range blogs {
		d, err := dto.NewInfoBlog(b)
		if err != nil {
			return nil, err
		}
		dtos[i] = d
	}

	return dtos, nil
}
