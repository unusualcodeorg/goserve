package micro

import (
	"github.com/nats-io/nats.go/micro"
	"github.com/unusualcodeorg/goserve/arch/network"
)

type NatsGroup = micro.Group
type NatsHandlerFunc = micro.HandlerFunc
type NatsRequest = micro.Request

type SendMessage interface {
	Message(data any)
	Error(err error)
}

type MessageSender interface {
	SendNats(req NatsRequest) SendMessage
}

type BaseController interface {
	MessageSender
	network.BaseController
}

type Controller interface {
	BaseController
	MountNats(group NatsGroup)
}

type Router interface {
	network.BaseRouter
	NatsClient() *NatsClient
	LoadControllers(controllers []Controller)
}

type Module[T any] interface {
	network.BaseModule[T]
	Controllers() []Controller
}
