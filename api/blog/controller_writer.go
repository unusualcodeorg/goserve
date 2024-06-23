package blog

import (
	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/goserve/api/blog/dto"
	userModel "github.com/unusualcodeorg/goserve/api/user/model"
	coredto "github.com/unusualcodeorg/goserve/arch/dto"
	"github.com/unusualcodeorg/goserve/arch/network"
	"github.com/unusualcodeorg/goserve/common"
)

type writerController struct {
	network.BaseController
	common.ContextPayload
	service Service
}

func NewWriterController(
	authMFunc network.AuthenticationProvider,
	authorizeMFunc network.AuthorizationProvider,
	service Service,
) network.Controller {
	return &writerController{
		BaseController: network.NewBaseController("/blog/writer", authMFunc, authorizeMFunc),
		ContextPayload: common.NewContextPayload(),
		service:        service,
	}
}

func (c *writerController) MountRoutes(group *gin.RouterGroup) {
	group.Use(c.Authentication(), c.Authorization(string(userModel.RoleCodeWriter)))
	group.POST("/", c.postBlogHandler)
	group.PUT("/", c.updateBlogHandler)
	group.GET("/id/:id", c.getBlogHandler)
	group.DELETE("/id/:id", c.deleteBlogHandler)
	group.PUT("/submit/id/:id", c.submitBlogHandler)
	group.PUT("/withdraw/id/:id", c.withdrawBlogHandler)
	group.GET("/drafts", c.getDraftsBlogsHandler)
	group.GET("/submitted", c.getSubmittedBlogsHandler)
	group.GET("/published", c.getPublishedBlogsHandler)
}

func (c *writerController) postBlogHandler(ctx *gin.Context) {
	body, err := network.ReqBody(ctx, dto.EmptyCreateBlog())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	user := c.MustGetUser(ctx)

	b, err := c.service.CreateBlog(body, user)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	c.Send(ctx).SuccessDataResponse("blog created successfully", b)
}

func (c *writerController) updateBlogHandler(ctx *gin.Context) {
	body, err := network.ReqBody(ctx, dto.EmptyUpdateBlog())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	user := c.MustGetUser(ctx)

	b, err := c.service.UpdateBlog(body, user)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	c.Send(ctx).SuccessDataResponse("blog updated successfully", b)
}

func (c *writerController) getBlogHandler(ctx *gin.Context) {
	mongoId, err := network.ReqParams(ctx, coredto.EmptyMongoId())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	user := c.MustGetUser(ctx)

	blog, err := c.service.GetBlogByIdForAuthor(mongoId.ID, user)
	if err != nil {
		c.Send(ctx).NotFoundError(mongoId.Id+" not found", err)
		return
	}

	c.Send(ctx).SuccessDataResponse("success", blog)
}

func (c *writerController) submitBlogHandler(ctx *gin.Context) {
	mongoId, err := network.ReqParams(ctx, coredto.EmptyMongoId())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	user := c.MustGetUser(ctx)

	err = c.service.BlogSubmission(mongoId.ID, user, true)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	c.Send(ctx).SuccessMsgResponse("blog submitted successfully")
}

func (c *writerController) withdrawBlogHandler(ctx *gin.Context) {
	mongoId, err := network.ReqParams(ctx, coredto.EmptyMongoId())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	user := c.MustGetUser(ctx)

	err = c.service.BlogSubmission(mongoId.ID, user, false)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	c.Send(ctx).SuccessMsgResponse("blog withdrawn successfully")
}

func (c *writerController) deleteBlogHandler(ctx *gin.Context) {
	mongoId, err := network.ReqParams(ctx, coredto.EmptyMongoId())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	user := c.MustGetUser(ctx)

	err = c.service.DeactivateBlog(mongoId.ID, user)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	c.Send(ctx).SuccessMsgResponse("blog deleted successfully")
}

func (c *writerController) getDraftsBlogsHandler(ctx *gin.Context) {
	pagination, err := network.ReqQuery(ctx, coredto.EmptyPagination())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	user := c.MustGetUser(ctx)

	blog, err := c.service.GetPaginatedDraftsForAuthor(user, pagination)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	c.Send(ctx).SuccessDataResponse("success", blog)
}

func (c *writerController) getSubmittedBlogsHandler(ctx *gin.Context) {
	pagination, err := network.ReqQuery(ctx, coredto.EmptyPagination())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	user := c.MustGetUser(ctx)

	blog, err := c.service.GetPaginatedSubmittedForAuthor(user, pagination)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	c.Send(ctx).SuccessDataResponse("success", blog)
}

func (c *writerController) getPublishedBlogsHandler(ctx *gin.Context) {
	pagination, err := network.ReqQuery(ctx, coredto.EmptyPagination())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	user := c.MustGetUser(ctx)

	blogs, err := c.service.GetPaginatedPublishedForAuthor(user, pagination)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	c.Send(ctx).SuccessDataResponse("success", blogs)
}
