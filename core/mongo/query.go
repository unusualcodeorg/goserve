package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseQuery[T any] interface {
	CreateIndexes(ctx context.Context, indexes []mongo.IndexModel) error
	FindOne(ctx context.Context, filter bson.M) (*T, error)
	FindAll(ctx context.Context, filter bson.M) ([]T, error)
	FindPaginated(ctx context.Context, filter bson.M, page int64, limit int64) ([]T, error)
	InsertOne(ctx context.Context, doc *T) (*primitive.ObjectID, error)
	InsertAndRetrieveOne(ctx context.Context, doc *T) (*T, error)
	InsertMany(ctx context.Context, doc []T) ([]primitive.ObjectID, error)
	InsertAndRetrieveMany(ctx context.Context, doc []T) ([]T, error)
	UpdateOne(ctx context.Context, filter bson.M, update bson.M) (int64, error)
	UpdateMany(ctx context.Context, filter bson.M, update bson.M) (int64, error)
	DeleteOne(ctx context.Context, filter bson.M) (int64, error)
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

func (q *query[T]) FindOne(ctx context.Context, filter bson.M) (*T, error) {
	collection := q.db.Collection(q.collectionName)

	var doc T
	err := collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		return nil, err
	}

	return &doc, nil
}

func (q *query[T]) FindAll(ctx context.Context, filter bson.M) ([]T, error) {
	collection := q.db.Collection(q.collectionName)

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []T

	for cursor.Next(ctx) {
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

	return docs, nil
}

func (q *query[T]) FindPaginated(ctx context.Context, filter bson.M, page int64, limit int64) ([]T, error) {
	collection := q.db.Collection(q.collectionName)

	skip := (page - 1) * limit

	findOptions := options.Find()
	findOptions.SetSkip(skip)
	findOptions.SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []T

	for cursor.Next(ctx) {
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

	return docs, nil
}

func (q *query[T]) InsertOne(ctx context.Context, doc *T) (*primitive.ObjectID, error) {
	collection := q.db.Collection(q.collectionName)

	result, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("database query error for: %s", insertedID)
	}

	return &insertedID, nil
}

func (q *query[T]) InsertAndRetrieveOne(ctx context.Context, doc *T) (*T, error) {
	collection := q.db.Collection(q.collectionName)

	result, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": result.InsertedID}
	retrived, err := q.FindOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	return retrived, nil
}

func (q *query[T]) InsertMany(ctx context.Context, docs []T) ([]primitive.ObjectID, error) {
	collection := q.db.Collection(q.collectionName)

	var iDocs []interface{}
	for _, doc := range docs {
		iDocs = append(iDocs, doc)
	}

	result, err := collection.InsertMany(ctx, iDocs)
	if err != nil {
		return nil, err
	}

	var insertedIDs []primitive.ObjectID

	for _, v := range result.InsertedIDs {
		insertedID, ok := v.(primitive.ObjectID)
		if !ok {
			return nil, fmt.Errorf("database query error for: %s", insertedID)
		}
		insertedIDs = append(insertedIDs, insertedID)
	}

	return insertedIDs, nil
}

func (q *query[T]) InsertAndRetrieveMany(ctx context.Context, docs []T) ([]T, error) {
	collection := q.db.Collection(q.collectionName)

	var iDocs []interface{}
	for _, doc := range docs {
		iDocs = append(iDocs, doc)
	}

	result, err := collection.InsertMany(ctx, iDocs)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": bson.M{"$in": result.InsertedIDs}}

	retrieved, err := q.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	return retrieved, nil
}

/*
 * Example -> update := bson.M{"$set": bson.M{"field": "newValue"}}
 */
func (q *query[T]) UpdateOne(ctx context.Context, filter bson.M, update bson.M) (int64, error) {
	collection := q.db.Collection(q.collectionName)

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	return result.MatchedCount, nil
}

/*
 * Example -> update := bson.M{"$set": bson.M{"field": "newValue"}}
 */
func (q *query[T]) UpdateMany(ctx context.Context, filter bson.M, update bson.M) (int64, error) {
	collection := q.db.Collection(q.collectionName)

	result, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	return result.MatchedCount, nil
}

func (q *query[T]) DeleteOne(ctx context.Context, filter bson.M) (int64, error) {
	collection := q.db.Collection(q.collectionName)

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
