package profile

import (
	"time"

	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProfileService interface {
	FindUserPrivateProfile(user *schema.User) (*schema.User, error)
	FindUserPubicProfile(user *schema.User) (*schema.User, error)
}

type service struct {
	network.BaseService
	userQuery mongo.Query[schema.User]
	roleQuery mongo.Query[schema.Role]
}

func NewProfileService(db mongo.Database, dbQueryTimeout time.Duration) ProfileService {
	s := service{
		BaseService: network.NewBaseService(dbQueryTimeout),
		userQuery:   mongo.NewQuery[schema.User](db, schema.UserCollectionName),
		roleQuery:   mongo.NewQuery[schema.Role](db, schema.RolesCollectionName),
	}
	return &s
}

func (s *service) FindUserPrivateProfile(user *schema.User) (*schema.User, error) {
	ctx, cancel := s.Context()
	defer cancel()
	filter := bson.M{"_id": user.ID, "status": true}
	return s.userQuery.FindOne(ctx, filter, nil)
}

func (s *service) FindUserPubicProfile(user *schema.User) (*schema.User, error) {
	ctx, cancel := s.Context()
	defer cancel()
	filter := bson.M{"_id": user.ID, "status": true}
	projection := bson.D{{Key: "name", Value: 1}, {Key: "profilePicUrl", Value: 1}}
	opts := options.FindOne().SetProjection(projection)
	return s.userQuery.FindOne(ctx, filter, opts)
}
