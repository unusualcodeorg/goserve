package network

import "github.com/gin-gonic/gin"

type baseMiddleware struct {
	ApiResponseSender
}

func NewBaseMiddleware() BaseMiddleware {
	return &baseMiddleware{
		ApiResponseSender: NewApiResponseSender(),
	}
}

func (m *baseMiddleware) Debug() bool {
	return gin.Mode() == gin.DebugMode
}
