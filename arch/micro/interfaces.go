package micro

import (
	"github.com/nats-io/nats.go/micro"
	"github.com/unusualcodeorg/goserve/arch/network"
)

type NatsGroup = micro.Group
type NatsHandlerFunc = micro.HandlerFunc
type NatsRequest = micro.Request

type BaseController interface {
	network.BaseController
	Context() *Context
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

type Module[T any] interface {
	network.BaseModule[T]
	Controllers() []Controller
}
