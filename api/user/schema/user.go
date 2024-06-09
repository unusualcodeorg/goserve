package schema

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

const CollectionName = "users"

type User struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty"`
	Name          *string              `bson:"name,omitempty" validate:"omitempty,max=200"`
	Email         string               `bson:"email" validate:"required,email"`
	Password      *string              `bson:"password" validate:"required,min=6,max=100"`
	ProfilePicURL *string              `bson:"profilePicUrl,omitempty" validate:"omitempty,max=500"`
	Roles         []primitive.ObjectID `bson:"roles,omitempty" validate:"required"`
	Verified      bool                 `bson:"verified" validate:"required"`
	Status        bool                 `bson:"status" validate:"required"`
	CreatedAt     time.Time            `bson:"createdAt" validate:"required"`
	UpdatedAt     time.Time            `bson:"updatedAt" validate:"required"`
}

func NewUser(email string, password *string) (*User, error) {
	now := time.Now()
	u := User{
		Email:     email,
		Password:  password,
		Verified:  false,
		Status:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := validateUser(u); err != nil {
		return nil, err
	}
	return &u, nil
}

func validateUser(u User) error {
	validate := validator.New()
	return validate.Struct(u)
}

func EnsureIndexes(db mongo.Database) {
	indexes := []mongod.IndexModel{
		{
			Keys: bson.D{
				{Key: "_id", Value: 1},
				{Key: "status", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "email", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	}
	q := mongo.NewDatabaseQuery[User](db, CollectionName)
	q.CreateIndexes(context.Background(), indexes)
}
