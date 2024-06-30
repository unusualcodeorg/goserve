package micro

import (
	"github.com/unusualcodeorg/goserve/arch/network"
)

type baseController struct {
	network.BaseController
	context *Context
}

func NewBaseController(basePath string, authProvider network.AuthenticationProvider, authorizeProvider network.AuthorizationProvider) BaseController {
	return &baseController{
		BaseController: network.NewBaseController(basePath, authProvider, authorizeProvider),
		context:        EmptyContext(),
	}
}

func (c *baseController) Context() *Context {
	return c.context
}
