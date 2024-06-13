package blog

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/blog/dto"
	coredto "github.com/unusualcodeorg/go-lang-backend-architecture/core/dto"
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
		blogService:  service,
	}
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.GET("/id/:id", c.getBlogHandler)
}

func (c *controller) getBlogHandler(ctx *gin.Context) {
	mongoId, err := network.ReqParams(ctx, &coredto.MongoId{})
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	blog, err := c.blogService.FindBlog(mongoId.ID)
	if err != nil {
		c.Send(ctx).NotFoundError("blog not found", err)
		return
	}

	data, err := network.MapToDto[dto.InfoBlog](blog)
	if err != nil {
		c.Send(ctx).InternalServerError("something went wrong", err)
		return
	}

	c.Send(ctx).SuccessDataResponse("success", data)
}
