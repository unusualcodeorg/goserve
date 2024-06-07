package network

type baseController struct {
	path                string
	authenticationMiddleware GroupMiddleware
	authorizationMiddleware  GroupMiddleware
}

func NewBaseController(
	path string,
	auth GroupMiddleware,
	authorize GroupMiddleware,
) BaseController {
	c := baseController{
		path:                path,
		authenticationMiddleware: auth,
		authorizationMiddleware:  authorize,
	}
	return &c
}

func (c *baseController) Path() string {
	return c.path
}

func (c *baseController) AuthenticationMiddleware() GroupMiddleware {
	return c.authenticationMiddleware
}

func (c *baseController) AuthorizationMiddleware() GroupMiddleware {
	return c.authorizationMiddleware
}
