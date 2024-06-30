package micro

import (
	"context"
)

type Context struct {
	context.Context
	NatsClient  *NatsClient
	NatsSubject string
}

func EmptyContext() *Context {
	return &Context{
		Context: context.Background(),
	}
}

func NewContext(natsClient *NatsClient, baseSub string) *Context {
	return &Context{
		Context:     context.Background(),
		NatsClient:  natsClient,
		NatsSubject: baseSub,
	}
}
