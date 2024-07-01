package micro

import (
	"github.com/unusualcodeorg/goserve/arch/network"
)

type baseController struct {
	MessageSender
	network.BaseController
}

func NewBaseController(basePath string, authProvider network.AuthenticationProvider, authorizeProvider network.AuthorizationProvider) BaseController {
	return &baseController{
		MessageSender:  NewMessageSender(),
		BaseController: network.NewBaseController(basePath, authProvider, authorizeProvider),
	}
}
