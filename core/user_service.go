package core

import (
	"time"

	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService interface {
	FindRoleByCode(code schema.RoleCode) (*schema.Role, error)
	FindRoles(roleIds []primitive.ObjectID) ([]schema.Role, error)
	FindUserById(id primitive.ObjectID) (*schema.User, error)
	FindUserByEmail(email string) (*schema.User, error)
	CreateUser(user *schema.User) (*schema.User, error)
	// UpdateUserInfo(user *schema.User) (*schema.User, error)
	// DeactivateUser(user *schema.User) (*schema.User, error)
}

type userService struct {
	network.BaseService
	userQuery mongo.Query[schema.User]
	roleQuery mongo.Query[schema.Role]
}

func NewUserService(db mongo.Database, dbQueryTimeout time.Duration) UserService {
	s := userService{
		BaseService: network.NewBaseService(dbQueryTimeout),
		userQuery:   mongo.NewQuery[schema.User](db, schema.UserCollectionName),
		roleQuery:   mongo.NewQuery[schema.Role](db, schema.RolesCollectionName),
	}
	return &s
}

func (s *userService) FindRoleByCode(code schema.RoleCode) (*schema.Role, error) {
	ctx, cancel := s.Context()
	defer cancel()
	filter := bson.M{"code": code, "status": true}
	return s.roleQuery.FindOne(ctx, filter, nil)
}

func (s *userService) FindRoles(roleIds []primitive.ObjectID) ([]schema.Role, error) {
	ctx, cancel := s.Context()
	cancel()
	filter := bson.M{"_id": bson.M{"$in": roleIds}}
	return s.roleQuery.FindAll(ctx, filter, nil)
}

func (s *userService) FindUserById(id primitive.ObjectID) (*schema.User, error) {
	ctx, cancel := s.Context()
	defer cancel()

	userFilter := bson.M{"_id": id, "status": true}
	user, err := s.userQuery.FindOne(ctx, userFilter, nil)
	if err != nil {
		return nil, err
	}

	roles, err := s.FindRoles(user.Roles)
	if err != nil {
		return nil, err
	}

	user.RoleDocs = roles
	return user, nil
}

func (s *userService) FindUserByEmail(email string) (*schema.User, error) {
	ctx, cancel := s.Context()
	defer cancel()
	filter := bson.M{"email": email, "status": true}
	return s.userQuery.FindOne(ctx, filter, nil)
}

func (s *userService) CreateUser(user *schema.User) (*schema.User, error) {
	ctx, cancel := s.Context()
	defer cancel()
	id, err := s.userQuery.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = *id
	return user, nil
}
