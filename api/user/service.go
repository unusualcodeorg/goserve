package user

import (
	"time"

	"github.com/unusualcodeorg/go-lang-backend-architecture/api/auth"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/user/schema"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserService interface {
	FindUserById(id primitive.ObjectID) (*schema.User, error)
	FindUserByemail(email string) (*schema.User, error)
	FindUserPrivateProfile(user *schema.User) (*schema.User, error)
	FindUserPubicProfile(user *schema.User) (*schema.User, error)
	// UpdateUserInfo(user *schema.User) (*schema.User, error)
	// DeactivateUser(user *schema.User) (*schema.User, error)
}

type service struct {
	network.BaseService
	authService auth.AuthService
	userQuery   mongo.Query[schema.User]
}

func NewUserService(db mongo.Database, dbQueryTimeout time.Duration, authService auth.AuthService) UserService {
	s := service{
		BaseService: network.NewBaseService(dbQueryTimeout),
		authService: authService,
		userQuery:   mongo.NewQuery[schema.User](db, schema.CollectionName),
	}
	return &s
}

func (s *service) FindUserById(id primitive.ObjectID) (*schema.User, error) {
	ctx, cancel := s.Context()
	defer cancel()

	userFilter := bson.M{"_id": id, "status": true}
	user, err := s.userQuery.FindOne(ctx, userFilter, nil)
	if err != nil {
		return nil, err
	}

	roles, err := s.authService.FindRoles(user.Roles)
	if err != nil {
		return nil, err
	}

	user.RoleDocs = roles
	return user, nil
}

func (s *service) FindUserByemail(email string) (*schema.User, error) {
	ctx, cancel := s.Context()
	defer cancel()
	filter := bson.M{"email": email, "status": true}
	return s.userQuery.FindOne(ctx, filter, nil)
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
