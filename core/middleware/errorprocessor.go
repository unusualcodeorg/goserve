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

func NewErrorProcessor() network.RootMiddleware {
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
					network.ResBadRequest(ctx, apiError.Message)
				case http.StatusForbidden:
					network.ResForbidden(ctx, apiError.Message)
				case http.StatusUnauthorized:
					network.ResUnauthorized(ctx, apiError.Message)
				case http.StatusNotFound:
					network.ResNotFound(ctx, apiError.Message)
				case http.StatusInternalServerError:
					if m.debug {
						network.ResInternalServerError(ctx, apiError.Message)
					} else {
						network.ResInternalServerError(ctx, "An unexpected error occurred. Please try again later.")
					}
				default:
					if m.debug {
						network.ResInternalServerError(ctx, apiError.Message)
					} else {
						network.ResInternalServerError(ctx, "An unexpected error occurred. Please try again later.")
					}
				}
			} else {
				network.ResInternalServerError(ctx, "An unexpected error occurred. Please try again later.")
			}
			ctx.Abort()
		}
	}()
	ctx.Next()
}
