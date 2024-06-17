package blog

import (
	"github.com/gin-gonic/gin"
	userModel "github.com/unusualcodeorg/go-lang-backend-architecture/api/user/model"
	"github.com/unusualcodeorg/go-lang-backend-architecture/common"
	coredto "github.com/unusualcodeorg/go-lang-backend-architecture/framework/dto"
	"github.com/unusualcodeorg/go-lang-backend-architecture/framework/network"
)

type editorController struct {
	network.BaseController
	common.ContextPayload
	service Service
}

func NewEditorController(
	authMFunc network.AuthenticationProvider,
	authorizeMFunc network.AuthorizationProvider,
	service Service,
) network.Controller {
	return &editorController{
		BaseController: network.NewBaseController("/blog/editor", authMFunc, authorizeMFunc),
		ContextPayload: common.NewContextPayload(),
		service:        service,
	}
}

func (c *editorController) MountRoutes(group *gin.RouterGroup) {
	group.Use(c.Authentication(), c.Authorization(string(userModel.RoleCodeEditor)))
	group.PUT("/publish/id/:id", c.publishBlogHandler)
	group.PUT("/unpublish/id/:id", c.unpublishBlogHandler)

}

func (c *editorController) publishBlogHandler(ctx *gin.Context) {
 mongoId, err := network.ReqParams(ctx, coredto.EmptyMongoId())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	user := c.MustGetUser(ctx)

	err = c.service.BlogPublicationForEditor(mongoId.ID, user, true)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	c.Send(ctx).SuccessMsgResponse("blog published successfully")
}

func (c *editorController) unpublishBlogHandler(ctx *gin.Context) {
 mongoId, err := network.ReqParams(ctx, coredto.EmptyMongoId())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	user := c.MustGetUser(ctx)

	err = c.service.BlogPublicationForEditor(mongoId.ID, user, false)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	c.Send(ctx).SuccessMsgResponse("blog unpublished successfully")
}
