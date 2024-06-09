package schema

import (
	"context"
	"time"

	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongod "go.mongodb.org/mongo-driver/mongo"
)

const KeystoreCollectionName = "keystores"

type Keystore struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Client       primitive.ObjectID `bson:"client" validate:"required"`
	PrimaryKey   string             `bson:"primaryKey" validate:"required"`
	SecondaryKey string             `bson:"secondaryKey" validate:"required"`
	Status       bool               `bson:"status" validate:"-"`
	CreatedAt    time.Time          `bson:"createdAt" validate:"required"`
	UpdatedAt    time.Time          `bson:"updatedAt" validate:"required"`
}

func NewKeystore(clientID primitive.ObjectID, primaryKey string, secondaryKey string) (*Keystore, error) {
	now := time.Now()
	k := Keystore{
		Client:       clientID,
		PrimaryKey:   primaryKey,
		SecondaryKey: secondaryKey,
		Status:       true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := utils.Validate(k); err != nil {
		return nil, err
	}
	return &k, nil
}

func EnsureKeystoreIndexes(db mongo.Database) {
	indexes := []mongod.IndexModel{
		{
			Keys: bson.D{
				{Key: "client", Value: 1},
				{Key: "primaryKey", Value: 1},
				{Key: "secondaryKey", Value: 1},
				{Key: "status", Value: 1},
			},
		},
	}
	q := mongo.NewDatabaseQuery[Role](db, CollectionName)
	q.CreateIndexes(context.Background(), indexes)
}
