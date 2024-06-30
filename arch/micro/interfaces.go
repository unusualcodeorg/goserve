package micro

import (
	"github.com/nats-io/nats.go/micro"
	"github.com/unusualcodeorg/goserve/arch/network"
)

type NatsGroup = micro.Group
type NatsHandlerFunc = micro.HandlerFunc
type NatsRequest = micro.Request

type Config struct {
	NatsUrl            string
	NatsServiceName    string
	NatsServiceVersion string
}

type BaseController interface {
	network.BaseController
	SetNatsContext(ctx *NatsContext)
	NatsContext() *NatsContext
}

type Controller interface {
	BaseController
	MountNats(group NatsGroup)
}

type Router interface {
	network.BaseRouter
	NatsClient() *NatsClient
	Disconnect()
	LoadControllers(controllers []Controller)
}
