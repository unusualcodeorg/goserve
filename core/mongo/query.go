package mongo

import (
	"context"
	"fmt"

	"github.com/unusualcodeorg/go-lang-backend-architecture/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseQuery[T any] interface {
	CreateIndexes(context context.Context, indexes []mongo.IndexModel) error
	FindOne(context context.Context, filter bson.M) (*T, error)
	FindPaginated(context context.Context, filter bson.M, page int64, limit int64) (*[]T, error)
	InsertOne(context context.Context, doc *T) (*T, error)
	// InsertMany(collectionName string, docs []any) ([]*primitive.ObjectID, error)
	// UpdateOne( collectionName string, filter any, doc any) (int64, error)
	// DeleteOne( collectionName string, filter any, doc any) (int64, error)
}

type query[T any] struct {
	db             Database
	collectionName string
}

func NewDatabaseQuery[T any](db Database, collectionName string) DatabaseQuery[T] {
	return &query[T]{
		db:             db,
		collectionName: collectionName,
	}
}

func (q *query[T]) CreateIndexes(ctx context.Context, indexes []mongo.IndexModel) error {
	fmt.Println("database creating index for: " + q.collectionName)
	collection := q.db.Collection(q.collectionName)
	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}

func (q *query[T]) FindOne(context context.Context, filter bson.M) (*T, error) {
	collection := q.db.Collection(q.collectionName)

	var doc T
	err := collection.FindOne(context, filter).Decode(&doc)
	if err != nil {
		return nil, err
	}

	return &doc, nil
}

func (q *query[T]) FindPaginated(context context.Context, filter bson.M, page int64, limit int64) (*[]T, error) {
	collection := q.db.Collection(q.collectionName)

	skip := (page - 1) * limit

	findOptions := options.Find()
	findOptions.SetSkip(skip)
	findOptions.SetLimit(int64(limit))

	cursor, err := collection.Find(context, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("error executing query1: %w", err)
	}
	defer cursor.Close(context)

	var docs []T

	for cursor.Next(context) {
		var result T
		err := cursor.Decode(&result)
		if err != nil {
			return nil, fmt.Errorf("error decoding result: %w", err)
		}
		docs = append(docs, result)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return &docs, nil
}

func (q *query[T]) InsertOne(context context.Context, doc *T) (*T, error) {
	collection := q.db.Collection(q.collectionName)

	result, err := collection.InsertOne(context, doc)
	if err != nil {
		return nil, err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("castObjectID: cast failed id")
	}

	d := utils.CopyAndSetField(doc, "ID", &insertedID)

	return d, nil
}

// func (q *query1) InsertMany( collectionName string, docs []any) ([]*primitive.ObjectID, error) {
// 	return nil, nil
// }

// func (q *query1) UpdateOne( collectionName string, filter any, doc any) (int64, error) {
// 	return 0, nil
// }

// func (q *query1) DeleteOne( collectionName string, filter any, doc any) (int64, error) {
// 	return 0, nil
// }
