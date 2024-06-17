package blog

import (
	"github.com/gin-gonic/gin"
	coredto "github.com/unusualcodeorg/goserve/framework/dto"
	"github.com/unusualcodeorg/goserve/framework/network"
)

type controller struct {
	network.BaseController
	service Service
}

func NewController(
	authMFunc network.AuthenticationProvider,
	authorizeMFunc network.AuthorizationProvider,
	service Service,
) network.Controller {
	return &controller{
		BaseController: network.NewBaseController("/blog", authMFunc, authorizeMFunc),
		service:        service,
	}
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.GET("/id/:id", c.getBlogByIdHandler)
	group.GET("/slug/:slug", c.getBlogBySlugHandler)

}

func (c *controller) getBlogByIdHandler(ctx *gin.Context) {
	mongoId, err := network.ReqParams(ctx, coredto.EmptyMongoId())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	blog, err := c.service.GetPublisedBlogById(mongoId.ID)
	if err != nil {
		c.Send(ctx).NotFoundError(mongoId.Id+" not found", err)
		return
	}

	c.Send(ctx).SuccessDataResponse("success", blog)
}

func (c *controller) getBlogBySlugHandler(ctx *gin.Context) {
	slug, err := network.ReqParams(ctx, coredto.EmptySlug())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	blog, err := c.service.GetPublishedBlogBySlug(slug.Slug)
	if err != nil {
		c.Send(ctx).NotFoundError(slug.Slug+" not found", err)
		return
	}

	c.Send(ctx).SuccessDataResponse("success", blog)
}
