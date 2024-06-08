package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DatabaseQuery interface {
	FindOne(ctx context.Context, collectionName string, filter any, doc any) error
	// FindAll(ctx context.Context, collectionName string, filter any) (any, error)
	InsertOne(ctx context.Context, collectionName string, doc any) (*primitive.ObjectID, error)
	// InsertMany(ctx context.Context, collectionName string, docs []any) ([]*primitive.ObjectID, error)
	// UpdateOne(ctx context.Context, collectionName string, filter any, doc any) (int64, error)
	// DeleteOne(ctx context.Context, collectionName string, filter any, doc any) (int64, error)
}

type query struct {
	db Database
}

func NewDatabaseQuery(db Database) DatabaseQuery {
	return &query{db: db}
}

func (q *query) FindOne(ctx context.Context, collectionName string, filter any, doc any) error {
	collection := q.db.GetCollection(collectionName)

	err := collection.FindOne(ctx, filter).Decode(doc)
	if err != nil {
		return err
	}

	return nil
}

// func (q *query) FindAll(ctx context.Context, collectionName string, filter any) ([]any, error) {
// 	return nil, nil
// }

func (q *query) InsertOne(ctx context.Context, collectionName string, doc any) (*primitive.ObjectID, error) {
	collection := q.db.GetCollection(collectionName)

	result, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, err
	}

	return &insertedID, nil
}

// func (q *query) InsertMany(ctx context.Context, collectionName string, docs []any) ([]*primitive.ObjectID, error) {
// 	return nil, nil
// }

// func (q *query) UpdateOne(ctx context.Context, collectionName string, filter any, doc any) (int64, error) {
// 	return 0, nil
// }

// func (q *query) DeleteOne(ctx context.Context, collectionName string, filter any, doc any) (int64, error) {
// 	return 0, nil
// }
