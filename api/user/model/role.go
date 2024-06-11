package model

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongod "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RoleCode string

const (
	RoleCodeLearner RoleCode = "LEARNER"
	RoleCodeAdmin   RoleCode = "ADMIN"
	RoleCodeWriter  RoleCode = "WRITER"
	RoleCodeEditor  RoleCode = "EDITOR"
)

type Role struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Code      RoleCode           `bson:"code" validate:"required,rolecode"`
	Status    bool               `bson:"status" validate:"required"`
	CreatedAt time.Time          `bson:"createdAt" validate:"required"`
	UpdatedAt time.Time          `bson:"updatedAt" validate:"required"`
}

const RolesCollectionName = "roles"

func NewRole(code RoleCode) (*Role, error) {
	now := time.Now()
	r := Role{
		Code:      code,
		Status:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := validateRole(r); err != nil {
		return nil, err
	}
	return &r, nil
}

func validateRole(r Role) error {
	validate := validator.New()

	_ = validate.RegisterValidation("rolecode", func(fl validator.FieldLevel) bool {
		code := RoleCode(fl.Field().String())
		switch code {
		case RoleCodeLearner, RoleCodeAdmin, RoleCodeWriter, RoleCodeEditor:
			return true
		}
		return false
	})

	return validate.Struct(r)
}

func (role *Role) GetValue() *Role {
	return role
}

func (role *Role) Validate() error {
	validate := validator.New()
	return validate.Struct(role)
}

func (*Role) EnsureIndexes(db mongo.Database) {
	indexes := []mongod.IndexModel{
		{
			Keys: bson.D{
				{Key: "_id", Value: 1},
				{Key: "status", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "code", Value: 1},
				{Key: "status", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "code", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	}
	q := mongo.NewQuery[Role](db, RolesCollectionName)
	q.CreateIndexes(context.Background(), indexes)
}
