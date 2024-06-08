package network

type baseController struct {
	path               string
	authenticationFunc GroupMiddlewareFunc
	authorizationFunc  GroupMiddlewareFunc
}

func NewBaseController(path string, authMFunc GroupMiddlewareFunc, authorizeMFunc GroupMiddlewareFunc) BaseController {
	c := baseController{
		path:               path,
		authenticationFunc: authMFunc,
		authorizationFunc:  authorizeMFunc,
	}
	return &c
}

func (c *baseController) Path() string {
	return c.path
}

func (c *baseController) AuthenticationMiddleware() GroupMiddleware {
	return c.authenticationFunc()
}

func (c *baseController) AuthorizationMiddleware() GroupMiddleware {
	return c.authorizationFunc()
}
