package contact

import (
	"context"

	"github.com/unusualcodeorg/go-lang-backend-architecture/api/contact/schema"
	"github.com/unusualcodeorg/go-lang-backend-architecture/internal/core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func saveMessage(msgType string, msgTxt string) (*schema.Message, error) {

	msg := schema.NewMessage(msgType, msgTxt)

	if err := core.Validate(msg); err != nil {
		return nil, err
	}

	collection := core.MongoCollection(msg.CollectionName)

	result, err := collection.InsertOne(context.Background(), msg)
	if err != nil {
		return nil, err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, err
	}

	msg.ID = insertedID.Hex()
	return msg, nil
}
