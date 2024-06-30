package micro

import (
	"context"
)

type NatsContext struct {
	context.Context
	Client  *NatsClient
	Subject string
}

func NewNatsContext(client *NatsClient, baseSub string) *NatsContext {
	return &NatsContext{
		Context: context.Background(),
		Client:  client,
		Subject: baseSub,
	}
}
