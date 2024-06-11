package network

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/schema"
)

type baseController struct {
	basePath          string
	authProvider      AuthenticationProvider
	authorizeProvider AuthorizationProvider
}

func NewBaseController(basePath string, authProvider AuthenticationProvider, authorizeProvider AuthorizationProvider) BaseController {
	c := baseController{
		basePath:          basePath,
		authProvider:      authProvider,
		authorizeProvider: authorizeProvider,
	}
	return &c
}

func (c *baseController) Path() string {
	return c.basePath
}

func (c *baseController) Authentication() gin.HandlerFunc {
	return c.authProvider.Middleware()
}

func (c *baseController) Authorization(roleCode schema.RoleCode) gin.HandlerFunc {
	return c.authorizeProvider.Middleware(roleCode)
}
