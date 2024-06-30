package micro

import (
	"github.com/unusualcodeorg/goserve/arch/network"
)

type baseController struct {
	network.BaseController
	natsCtx *NatsContext
}

func NewBaseController(basePath string, authProvider network.AuthenticationProvider, authorizeProvider network.AuthorizationProvider) BaseController {
	return &baseController{
		BaseController: network.NewBaseController(basePath, authProvider, authorizeProvider),
	}
}

func (c *baseController) SetNatsContext(ctx *NatsContext) {
	c.natsCtx = ctx
}

func (c *baseController) NatsContext() *NatsContext {
	return c.natsCtx
}
