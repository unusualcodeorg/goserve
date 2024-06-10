package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

type errorHandle struct {
	debug bool
}

func NewErrorHandle() network.RootMiddleware {
	debug := gin.Mode() == gin.DebugMode
	m := errorHandle{debug: debug}
	return &m
}

func (m *errorHandle) Attach(engine *gin.Engine) {
	engine.Use(m.Handler)
}

func (m *errorHandle) Handler(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			var apiError *network.ApiError
			if errors.As(r.(error), &apiError) {
				switch apiError.Code {
				case http.StatusBadRequest:
					network.BadRequestResponse(apiError.Message).Send(ctx)
				case http.StatusForbidden:
					network.ForbiddenResponse(apiError.Message).Send(ctx)
				case http.StatusUnauthorized:
					network.UnauthorizedResponse(apiError.Message).Send(ctx)
				case http.StatusNotFound:
					network.NotFoundResponse(apiError.Message).Send(ctx)
				case http.StatusInternalServerError:
					if m.debug {
						network.InternalServerErrorResponse(apiError.Message).Send(ctx)
					} else {
						network.InternalServerErrorResponse("An unexpected error occurred. Please try again later.").Send(ctx)
					}
				default:
					if m.debug {
						network.InternalServerErrorResponse(apiError.Message).Send(ctx)
					} else {
						network.InternalServerErrorResponse("An unexpected error occurred. Please try again later.").Send(ctx)
					}
				}
			} else {
				network.InternalServerErrorResponse("An unexpected error occurred. Please try again later.").Send(ctx)
			}
			ctx.Abort()
		}
	}()
	ctx.Next()
}
