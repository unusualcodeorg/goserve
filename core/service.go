package core

import "github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"

type CoreService interface {
}

type service struct {
	db mongo.Database
}

func NewCoreService(database mongo.Database) CoreService {
	s := service{
		db: database,
	}
	return &s
}
