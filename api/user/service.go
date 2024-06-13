package user

import (
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/user/model"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
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

func NewService(db mongo.Database) Service {
	s := service{
		BaseService:      network.NewBaseService(),
		userQueryBuilder: mongo.NewQueryBuilder[model.User](db, model.UserCollectionName),
		roleQueryBuilder: mongo.NewQueryBuilder[model.Role](db, model.RolesCollectionName),
	}
	return &s
}

func (s *service) FindRoleByCode(code model.RoleCode) (*model.Role, error) {
	filter := bson.M{"code": code, "status": true}
	return s.roleQueryBuilder.SingleQuery().FindOne(filter, nil)
}

func (s *service) FindRoles(roleIds []primitive.ObjectID) ([]model.Role, error) {
	filter := bson.M{"_id": bson.M{"$in": roleIds}}
	return s.roleQueryBuilder.SingleQuery().FindAll(filter, nil)
}

func (s *service) FindUserById(id primitive.ObjectID) (*model.User, error) {
	userFilter := bson.M{"_id": id, "status": true}
	proj := bson.D{{Key: "password", Value: 0}}
	opts := options.FindOne().SetProjection(proj)
	user, err := s.userQueryBuilder.SingleQuery().FindOne(userFilter, opts)
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
	filter := bson.M{"email": email, "status": true}
	user, err := s.userQueryBuilder.SingleQuery().FindOne(filter, nil)

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
	id, err := s.userQueryBuilder.SingleQuery().InsertOne(user)
	if err != nil {
		return nil, err
	}
	user.ID = *id
	return user, nil
}

func (s *service) FindUserPrivateProfile(user *model.User) (*model.User, error) {
	filter := bson.M{"_id": user.ID, "status": true}
	projection := bson.D{{Key: "password", Value: 0}}
	opts := options.FindOne().SetProjection(projection)
	return s.userQueryBuilder.SingleQuery().FindOne(filter, opts)
}

func (s *service) FindUserPubicProfile(user *model.User) (*model.User, error) {
	filter := bson.M{"_id": user.ID, "status": true}
	projection := bson.D{{Key: "name", Value: 1}, {Key: "profilePicUrl", Value: 1}}
	opts := options.FindOne().SetProjection(projection)
	return s.userQueryBuilder.SingleQuery().FindOne(filter, opts)
}
