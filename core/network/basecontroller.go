package network

type baseController struct {
	path               string
	authenticationFunc GroupMiddlewareFunc
	authorizationFunc  GroupMiddlewareFunc
}

func NewBaseController(
	path string,
	authenticationFunc GroupMiddlewareFunc,
	authorizationFunc GroupMiddlewareFunc,
) BaseController {
	c := baseController{
		path:               path,
		authenticationFunc: authenticationFunc,
		authorizationFunc:  authorizationFunc,
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
