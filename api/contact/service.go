package contact

import (
	"context"
	"time"

	"github.com/unusualcodeorg/go-lang-backend-architecture/api/contact/schema"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
)

type ContactService interface {
	SaveMessage(msgType string, msgTxt string) (*schema.Message, error)
}

type service struct {
	dbQuery mongo.DatabaseQuery
}

func NewContactService(dbQuery mongo.DatabaseQuery) ContactService {
	s := service{
		dbQuery: dbQuery,
	}
	return &s
}

func (s *service) SaveMessage(msgType string, msgTxt string) (*schema.Message, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	msg, err := schema.NewMessage(msgType, msgTxt)
	if err != nil {
		return nil, err
	}

	result, err := s.dbQuery.InsertOne(ctx, schema.MessageCollectionName, msg)
	if err != nil {
		return nil, err
	}

	msg.ID = result.Hex()
	return msg, nil
}
