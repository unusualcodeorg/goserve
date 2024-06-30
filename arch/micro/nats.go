package micro

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
)

type Config struct {
	NatsUrl            string
	NatsServiceName    string
	NatsServiceVersion string
}

type NatsClient struct {
	Conn        *nats.Conn
	Service     micro.Service
	Timeout time.Duration
}

func NewNatsClient(config *Config) *NatsClient {
	fmt.Println("connecting to nats..")

	nc, err := nats.Connect(config.NatsUrl)
	if err != nil {
		panic(err)
	}

	srv, err := micro.AddService(nc, micro.Config{
		Name:    config.NatsServiceName,
		Version: config.NatsServiceVersion,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("connected to nats..")

	return &NatsClient{
		Conn:        nc,
		Service:     srv,
		Timeout: nats.DefaultTimeout,
	}
}
