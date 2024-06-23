package blog

import (
	"time"

	"github.com/unusualcodeorg/goserve/api/blog/dto"
	"github.com/unusualcodeorg/goserve/api/blog/model"
	"github.com/unusualcodeorg/goserve/api/user"
	userModel "github.com/unusualcodeorg/goserve/api/user/model"
	coredto "github.com/unusualcodeorg/goserve/arch/dto"
	"github.com/unusualcodeorg/goserve/arch/mongo"
	"github.com/unusualcodeorg/goserve/arch/network"
	"github.com/unusualcodeorg/goserve/arch/redis"
	"github.com/unusualcodeorg/goserve/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	SetBlogDtoCacheById(blog *dto.PublicBlog) error
	GetBlogDtoCacheById(id primitive.ObjectID) (*dto.PublicBlog, error)
	SetBlogDtoCacheBySlug(blog *dto.PublicBlog) error
	GetBlogDtoCacheBySlug(slug string) (*dto.PublicBlog, error)
	SetSimilarBlogsDtoCache(blogId primitive.ObjectID, blogs []*dto.InfoBlog) error
	GetSimilarBlogsDtoCache(blogId primitive.ObjectID) ([]*dto.InfoBlog, error)
	BlogSlugExists(slug string) bool
	CreateBlog(createBlogDto *dto.CreateBlog, author *userModel.User) (*dto.PrivateBlog, error)
	UpdateBlog(updateBlogDto *dto.UpdateBlog, author *userModel.User) (*dto.PrivateBlog, error)
	DeactivateBlog(blogId primitive.ObjectID, author *userModel.User) error
	BlogSubmission(blogId primitive.ObjectID, author *userModel.User, submit bool) error
	GetPublisedBlogById(id primitive.ObjectID) (*dto.PublicBlog, error)
	GetPublishedBlogBySlug(slug string) (*dto.PublicBlog, error)
	GetBlogByIdForAuthor(id primitive.ObjectID, author *userModel.User) (*dto.PrivateBlog, error)
	GetPaginatedDraftsForAuthor(author *userModel.User, p *coredto.Pagination) ([]*dto.InfoBlog, error)
	GetPaginatedPublishedForAuthor(author *userModel.User, p *coredto.Pagination) ([]*dto.InfoBlog, error)
	GetPaginatedSubmittedForAuthor(author *userModel.User, p *coredto.Pagination) ([]*dto.InfoBlog, error)
	GetBlogByIdForEditor(id primitive.ObjectID) (*dto.PrivateBlog, error)
	BlogPublicationForEditor(blogId primitive.ObjectID, editor *userModel.User, publish bool) error
	GetPaginatedPublishedForEditor(p *coredto.Pagination) ([]*dto.InfoBlog, error)
	GetPaginatedSubmittedForEditor(p *coredto.Pagination) ([]*dto.InfoBlog, error)
	GetPaginatedLatestBlogs(p *coredto.Pagination) ([]*dto.InfoBlog, error)
	GetPaginatedTaggedBlogs(tag string, p *coredto.Pagination) ([]*dto.InfoBlog, error)
	GetSimilarBlogs(blogId primitive.ObjectID) ([]*dto.InfoBlog, error)
	getPublicPublishedBlog(filter bson.M) (*dto.PublicBlog, error)
	getPublicPaginated(filter bson.M, p *coredto.Pagination) ([]*dto.InfoBlog, error)
	getPaginated(filter bson.M, p *coredto.Pagination, opts *options.FindOptions) ([]*dto.InfoBlog, error)
}

type service struct {
	network.BaseService
	blogQueryBuilder mongo.QueryBuilder[model.Blog]
	publicBlogCache  redis.Cache[dto.PublicBlog]
	infoBlogCache    redis.Cache[dto.InfoBlog]
	userService      user.Service
}

func NewService(db mongo.Database, store redis.Store, userService user.Service) Service {
	return &service{
		BaseService:      network.NewBaseService(),
		blogQueryBuilder: mongo.NewQueryBuilder[model.Blog](db, model.CollectionName),
		publicBlogCache:  redis.NewCache[dto.PublicBlog](store),
		infoBlogCache:    redis.NewCache[dto.InfoBlog](store),
		userService:      userService,
	}
}

func (s *service) SetBlogDtoCacheById(blog *dto.PublicBlog) error {
	key := "blog_" + blog.ID.Hex()
	return s.publicBlogCache.SetJSON(key, blog, time.Duration(10*time.Minute))
}

func (s *service) GetBlogDtoCacheById(id primitive.ObjectID) (*dto.PublicBlog, error) {
	key := "blog_" + id.Hex()
	return s.publicBlogCache.GetJSON(key)
}

func (s *service) SetBlogDtoCacheBySlug(blog *dto.PublicBlog) error {
	key := "blog_" + blog.Slug
	return s.publicBlogCache.SetJSON(key, blog, time.Duration(10*time.Minute))
}

func (s *service) GetBlogDtoCacheBySlug(slug string) (*dto.PublicBlog, error) {
	key := "blog_" + slug
	return s.publicBlogCache.GetJSON(key)
}

func (s *service) SetSimilarBlogsDtoCache(blogId primitive.ObjectID, blogs []*dto.InfoBlog) error {
	key := "similar_blogs_" + blogId.Hex()
	return s.infoBlogCache.SetJSONList(key, blogs, 6*time.Hour)
}

func (s *service) GetSimilarBlogsDtoCache(blogId primitive.ObjectID) ([]*dto.InfoBlog, error) {
	key := "similar_blogs_" + blogId.Hex()
	return s.infoBlogCache.GetJSONList(key)
}

func (s *service) BlogSlugExists(slug string) bool {
	filter := bson.M{"slug": slug}
	projection := bson.D{{Key: "status", Value: 1}}
	opts := options.FindOne().SetProjection(projection)
	_, err := s.blogQueryBuilder.SingleQuery().FindOne(filter, opts)
	return err == nil
}

func (s *service) CreateBlog(b *dto.CreateBlog, author *userModel.User) (*dto.PrivateBlog, error) {
	b.Slug = utils.FormatEndpoint(b.Slug)

	exists := s.BlogSlugExists(b.Slug)
	if exists {
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

func (s *service) UpdateBlog(b *dto.UpdateBlog, author *userModel.User) (*dto.PrivateBlog, error) {
	filter := bson.M{"_id": b.ID, "author": author.ID, "status": true}
	blog, err := s.blogQueryBuilder.SingleQuery().FindOne(filter, nil)
	if err != nil {
		return nil, network.NewNotFoundError("Blog with id: "+b.ID.Hex()+" does not exists", nil)
	}

	updates := bson.M{}

	if b.Slug != nil {
		slug := utils.FormatEndpoint(*b.Slug)
		if slug != blog.Slug {
			exists := s.BlogSlugExists(slug)
			if exists {
				return nil, network.NewBadRequestError("Blog with slug: "+slug+" already exists", nil)
			}
			updates["slug"] = slug
		}
	}

	if b.Title != nil {
		updates["title"] = *b.Title
	}

	if b.Description != nil {
		updates["description"] = *b.Description
	}

	if b.DraftText != nil {
		updates["draftText"] = *b.DraftText
	}

	if b.Tags != nil {
		updates["tags"] = *b.Tags
	}

	if b.ImgURL != nil {
		updates["imgUrl"] = *b.ImgURL
	}

	updates["updatedBy"] = author.ID
	updates["updatedAt"] = time.Now()

	set := bson.M{"$set": updates}
	_, err = s.blogQueryBuilder.SingleQuery().UpdateOne(filter, set)
	if err != nil {
		return nil, err
	}

	return s.GetBlogByIdForAuthor(blog.ID, author)
}

func (s *service) DeactivateBlog(blogId primitive.ObjectID, author *userModel.User) error {
	filter := bson.M{"_id": blogId, "author": author.ID, "status": true}
	update := bson.M{"$set": bson.M{"status": false, "updatedBy": author.ID, "updatedAt": time.Now()}}
	result, err := s.blogQueryBuilder.SingleQuery().UpdateOne(filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return network.NewNotFoundError("blog not found", nil)
	}

	return nil
}

func (s *service) BlogSubmission(blogId primitive.ObjectID, author *userModel.User, submit bool) error {
	filter := bson.M{"_id": blogId, "author": author.ID, "status": true}
	update := bson.M{"$set": bson.M{"submitted": submit, "updatedBy": author.ID, "updatedAt": time.Now()}}
	result, err := s.blogQueryBuilder.SingleQuery().UpdateOne(filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return network.NewNotFoundError("blog not found", nil)
	}

	return nil
}

func (s *service) BlogPublicationForEditor(blogId primitive.ObjectID, editor *userModel.User, publish bool) error {
	filter := bson.M{"_id": blogId, "status": true}
	blog, err := s.blogQueryBuilder.SingleQuery().FindOne(filter, nil)
	if err != nil {
		return network.NewNotFoundError("blog for _id "+blogId.Hex()+" not found", err)
	}

	if !blog.Submitted {
		return network.NewBadRequestError("blog for _id "+blogId.Hex()+" is not submitted", err)
	}

	var update bson.M

	if publish {
		if blog.PublishedAt == nil {
			now := time.Now()
			blog.PublishedAt = &now
		}
		update = bson.M{"drafted": false, "submitted": false, "published": true, "text": blog.DraftText, "publishedAt": blog.PublishedAt}
	} else {
		update = bson.M{"drafted": true, "submitted": false, "published": false}
	}

	update["updatedBy"] = editor.ID
	update["updatedAt"] = time.Now()

	updated := bson.M{"$set": update}
	result, err := s.blogQueryBuilder.SingleQuery().UpdateOne(filter, updated)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return network.NewNotFoundError("blog not found", nil)
	}

	return nil
}

func (s *service) GetBlogByIdForAuthor(id primitive.ObjectID, author *userModel.User) (*dto.PrivateBlog, error) {
	filter := bson.M{"_id": id, "author": author.ID, "status": true}

	blog, err := s.blogQueryBuilder.SingleQuery().FindOne(filter, nil)
	if err != nil {
		return nil, err
	}

	return dto.NewPrivateBlog(blog, author)
}

func (s *service) GetPublisedBlogById(id primitive.ObjectID) (*dto.PublicBlog, error) {
	filter := bson.M{"_id": id, "published": true, "status": true}
	return s.getPublicPublishedBlog(filter)
}

func (s *service) GetPublishedBlogBySlug(slug string) (*dto.PublicBlog, error) {
	filter := bson.M{"slug": slug, "published": true, "status": true}
	return s.getPublicPublishedBlog(filter)
}

func (s *service) getPublicPublishedBlog(filter bson.M) (*dto.PublicBlog, error) {
	projection := bson.D{{Key: "draftText", Value: 0}}
	opts := options.FindOne().SetProjection(projection)
	blog, err := s.blogQueryBuilder.SingleQuery().FindOne(filter, opts)
	if err != nil {
		return nil, network.NewNotFoundError("blog not found", err)
	}

	author, err := s.userService.FindUserPublicProfile(blog.Author)
	if err != nil {
		return nil, network.NewNotFoundError("author not found", err)
	}

	return dto.NewPublicBlog(blog, author)
}

func (s *service) GetBlogByIdForEditor(id primitive.ObjectID) (*dto.PrivateBlog, error) {
	filter := bson.M{"_id": id, "status": true}
	blog, err := s.blogQueryBuilder.SingleQuery().FindOne(filter, nil)
	if err != nil {
		return nil, err
	}

	author, err := s.userService.FindUserPublicProfile(blog.Author)
	if err != nil {
		return nil, err
	}

	return dto.NewPrivateBlog(blog, author)
}

func (s *service) GetPaginatedDraftsForAuthor(author *userModel.User, p *coredto.Pagination) ([]*dto.InfoBlog, error) {
	filter := bson.M{"author": author.ID, "status": true, "drafted": true}
	return s.getPaginated(filter, p, nil)
}

func (s *service) GetPaginatedPublishedForAuthor(author *userModel.User, p *coredto.Pagination) ([]*dto.InfoBlog, error) {
	filter := bson.M{"author": author.ID, "status": true, "published": true}
	return s.getPaginated(filter, p, nil)
}

func (s *service) GetPaginatedSubmittedForAuthor(author *userModel.User, p *coredto.Pagination) ([]*dto.InfoBlog, error) {
	filter := bson.M{"author": author.ID, "status": true, "submitted": true}
	return s.getPaginated(filter, p, nil)
}

func (s *service) GetPaginatedPublishedForEditor(p *coredto.Pagination) ([]*dto.InfoBlog, error) {
	filter := bson.M{"status": true, "published": true}
	return s.getPaginated(filter, p, nil)
}

func (s *service) GetPaginatedSubmittedForEditor(p *coredto.Pagination) ([]*dto.InfoBlog, error) {
	filter := bson.M{"status": true, "submitted": true}
	return s.getPaginated(filter, p, nil)
}

func (s *service) GetPaginatedLatestBlogs(p *coredto.Pagination) ([]*dto.InfoBlog, error) {
	filter := bson.M{"status": true, "published": true}
	return s.getPublicPaginated(filter, p)
}

func (s *service) GetPaginatedTaggedBlogs(tag string, p *coredto.Pagination) ([]*dto.InfoBlog, error) {
	filter := bson.M{"status": true, "published": true, "tags": tag}
	return s.getPublicPaginated(filter, p)
}

func (s *service) GetSimilarBlogs(blogId primitive.ObjectID) ([]*dto.InfoBlog, error) {
	ftr := bson.M{"_id": blogId, "published": true, "status": true}
	dto, err := s.getPublicPublishedBlog(ftr)
	if err != nil {
		return nil, network.NewNotFoundError("blog not found", err)
	}

	filter := bson.M{
		"$text":     bson.M{"$search": dto.Title, "$caseSensitive": false},
		"status":    true,
		"published": true,
		"_id":       bson.M{"$ne": dto.ID},
	}

	opts := options.Find()
	opts.SetProjection(bson.M{"similarity": bson.M{"$meta": "textScore"}})
	opts.SetSort(bson.D{
		{Key: "similarity", Value: bson.M{"$meta": "textScore"}},
		{Key: "updatedAt", Value: -1},
		{Key: "score", Value: -1},
	})

	pagination := &coredto.Pagination{
		Page:  1,
		Limit: 6,
	}

	return s.getPaginated(filter, pagination, opts)
}

func (s *service) getPublicPaginated(filter bson.M, p *coredto.Pagination) ([]*dto.InfoBlog, error) {
	projection := bson.D{{Key: "draftText", Value: 0}}
	opts := options.Find().SetProjection(projection)
	opts.SetSort(bson.D{{Key: "updatedAt", Value: -1}, {Key: "score", Value: -1}})
	return s.getPaginated(filter, p, opts)
}

func (s *service) getPaginated(filter bson.M, p *coredto.Pagination, opts *options.FindOptions) ([]*dto.InfoBlog, error) {
	blogs, err := s.blogQueryBuilder.SingleQuery().FindPaginated(filter, p.Page, p.Limit, opts)
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
