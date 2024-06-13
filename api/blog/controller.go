package blog

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/blog/dto"
	userModel "github.com/unusualcodeorg/go-lang-backend-architecture/api/user/model"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

type controller struct {
	network.BaseController
	blogService BlogService
}

func NewBlogController(
	authMFunc network.AuthenticationProvider,
	authorizeMFunc network.AuthorizationProvider,
	service BlogService,
) network.Controller {
	return &controller{
		BaseController: network.NewBaseController("/blog", authMFunc, authorizeMFunc),
		blogService:    service,
	}
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	// group.GET("/id/:id", c.getBlogHandler)

	writer := group.Group("/writer", c.Authentication(), c.Authorization(string(userModel.RoleCodeWriter)))
	writer.POST("/", c.postBlogHandler)
}

func (c *controller) postBlogHandler(ctx *gin.Context) {
	body, err := network.ReqBody(ctx, dto.EmptyCreateBlog())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	user := network.ReqMustGetUser[userModel.User](ctx)

	b, err := c.blogService.CreateBlog(body, user)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	c.Send(ctx).SuccessDataResponse("blog creation success", b)
}
