package network

type baseController struct {
	path                 string
	authenticationRecipe GroupMiddlewareRecipe
	authorizationRecipe  GroupMiddlewareRecipe
}

func NewBaseController(
	path string,
	authRecipe GroupMiddlewareRecipe,
	authorizeRecipe GroupMiddlewareRecipe,
) BaseController {
	c := baseController{
		path:                 path,
		authenticationRecipe: authRecipe,
		authorizationRecipe:  authorizeRecipe,
	}
	return &c
}

func (c *baseController) Path() string {
	return c.path
}

func (c *baseController) AuthenticationMiddleware() GroupMiddleware {
	return c.authenticationRecipe()
}

func (c *baseController) AuthorizationMiddleware() GroupMiddleware {
	return c.authorizationRecipe()
}
