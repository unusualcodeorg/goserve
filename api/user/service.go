package user

import (
	"time"

	"github.com/unusualcodeorg/go-lang-backend-architecture/api/user/model"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserService interface {
	FindRoleByCode(code model.RoleCode) (*model.Role, error)
	FindRoles(roleIds []primitive.ObjectID) ([]model.Role, error)
	FindUserById(id primitive.ObjectID) (*model.User, error)
	FindUserByEmail(email string) (*model.User, error)
	CreateUser(user *model.User) (*model.User, error)
	FindUserPrivateProfile(user *model.User) (*model.User, error)
	FindUserPubicProfile(user *model.User) (*model.User, error)
	// UpdateUserInfo(user *model.User) (*model.User, error)
	// DeactivateUser(user *model.User) (*model.User, error)
}

type service struct {
	network.BaseService
	userQueryBuilder mongo.QueryBuilder[model.User]
	roleQueryBuilder mongo.QueryBuilder[model.Role]
}

func NewUserService(db mongo.Database, dbQueryTimeout time.Duration) UserService {
	s := service{
		BaseService:      network.NewBaseService(dbQueryTimeout),
		userQueryBuilder: mongo.NewQueryBuilder[model.User](db, model.UserCollectionName),
		roleQueryBuilder: mongo.NewQueryBuilder[model.Role](db, model.RolesCollectionName),
	}
	return &s
}

func (s *service) FindRoleByCode(code model.RoleCode) (*model.Role, error) {
	ctx, cancel := s.Context()
	defer cancel()
	filter := bson.M{"code": code, "status": true}
	return s.roleQueryBuilder.Query(ctx).FindOne(filter, nil)
}

func (s *service) FindRoles(roleIds []primitive.ObjectID) ([]model.Role, error) {
	ctx, cancel := s.Context()
	defer cancel()
	filter := bson.M{"_id": bson.M{"$in": roleIds}}
	return s.roleQueryBuilder.Query(ctx).FindAll(filter, nil)
}

func (s *service) FindUserById(id primitive.ObjectID) (*model.User, error) {
	ctx, cancel := s.Context()
	defer cancel()

	userFilter := bson.M{"_id": id, "status": true}
	proj := bson.D{{Key: "password", Value: 0}}
	opts := options.FindOne().SetProjection(proj)
	user, err := s.userQueryBuilder.Query(ctx).FindOne(userFilter, opts)
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

func (s *service) FindUserByEmail(email string) (*model.User, error) {
	ctx, cancel := s.Context()
	defer cancel()
	filter := bson.M{"email": email, "status": true}
	user, err := s.userQueryBuilder.Query(ctx).FindOne(filter, nil)

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

func (s *service) CreateUser(user *model.User) (*model.User, error) {
	ctx, cancel := s.Context()
	defer cancel()
	id, err := s.userQueryBuilder.Query(ctx).InsertOne(user)
	if err != nil {
		return nil, err
	}
	user.ID = *id
	return user, nil
}

func (s *service) FindUserPrivateProfile(user *model.User) (*model.User, error) {
	ctx, cancel := s.Context()
	defer cancel()
	filter := bson.M{"_id": user.ID, "status": true}
	projection := bson.D{{Key: "password", Value: 0}}
	opts := options.FindOne().SetProjection(projection)
	return s.userQueryBuilder.Query(ctx).FindOne(filter, opts)
}

func (s *service) FindUserPubicProfile(user *model.User) (*model.User, error) {
	ctx, cancel := s.Context()
	defer cancel()
	filter := bson.M{"_id": user.ID, "status": true}
	projection := bson.D{{Key: "name", Value: 1}, {Key: "profilePicUrl", Value: 1}}
	opts := options.FindOne().SetProjection(projection)
	return s.userQueryBuilder.Query(ctx).FindOne(filter, opts)
}
