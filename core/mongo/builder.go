package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type QueryBuilder[T any] interface {
	GetCollection() *mongo.Collection
	Query(context context.Context) Query[T]
}

type queryBuilder[T any] struct {
	collection *mongo.Collection
}

func (c *queryBuilder[T]) GetCollection() *mongo.Collection {
	return c.collection
}

func (c *queryBuilder[T]) Query(context context.Context) Query[T] {
	return newQuery[T](context, c.collection)
}

func NewQueryBuilder[T any](db Database, collectionName string) QueryBuilder[T] {
	return &queryBuilder[T]{
		collection: db.GetInstance().Collection(collectionName),
	}
}
