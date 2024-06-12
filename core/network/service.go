package network

import (
	"context"
	"time"
)

type baseService struct {
	dbQueryTimeout time.Duration
}

func NewBaseService(dbQueryTimeout time.Duration) BaseService {
	return &baseService{
		dbQueryTimeout: dbQueryTimeout,
	}
}

func (s *baseService) Context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), s.dbQueryTimeout)
}
