package network

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type apiResponseSender struct{}

func NewApiResponseSender() ApiResponseSender {
	return &apiResponseSender{}
}

func (m *apiResponseSender) Debug() bool {
	return gin.Mode() == gin.DebugMode
}

func (s *apiResponseSender) SendError(ctx *gin.Context, err ApiError) {
	var res ApiResponse
	var apiError ApiError

	if errors.As(err, &apiError) {
		e := apiError.GetValue()
		switch e.Code {
		case http.StatusBadRequest:
			res = BadRequestResponse(e.Message)
		case http.StatusForbidden:
			res = ForbiddenResponse(e.Message)
		case http.StatusUnauthorized:
			res = UnauthorizedResponse(e.Message)
		case http.StatusNotFound:
			res = NotFoundResponse(e.Message)
		case http.StatusInternalServerError:
			if s.Debug() {
				res = InternalServerErrorResponse(apiError.Unwrap().Error())
			}
		default:
			if s.Debug() {
				res = InternalServerErrorResponse(apiError.Unwrap().Error())
			}
		}
	}

	if res == nil {
		res = InternalServerErrorResponse("An unexpected error occurred. Please try again later.")
	}

	s.SendResponse(ctx, res)
}

func (s *apiResponseSender) SendResponse(ctx *gin.Context, response ApiResponse) {
	ctx.JSON(int(response.GetValue().Status), response)
}
