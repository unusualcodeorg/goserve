package network

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type sender struct{}

func NewResponseSender() ResponseSender {
	return &sender{}
}

func (m *sender) Debug() bool {
	return gin.Mode() == gin.DebugMode
}

func (m *sender) Send(ctx *gin.Context) ResponseSend {
	return &send{
		debug:   m.Debug(),
		context: ctx,
	}
}

type send struct {
	debug   bool
	context *gin.Context
}

func (s *send) SuccessMsgResponse(message string) {
	s.sendResponse(NewSuccessMsgResponse(message))
}

func (s *send) SuccessDataResponse(message string, data any) {
	s.sendResponse(NewSuccessDataResponse(message, data))
}

func (s *send) BadRequestError(message string, err error) {
	s.sendError(NewBadRequestError(message, err))
}

func (s *send) ForbiddenError(message string, err error) {
	s.sendError(NewForbiddenError(message, err))
}

func (s *send) UnauthorizedError(message string, err error) {
	s.sendError(NewUnauthorizedError(message, err))
}

func (s *send) NotFoundError(message string, err error) {
	s.sendError(NewNotFoundError(message, err))
}

func (s *send) InternalServerError(message string, err error) {
	s.sendError(NewInternalServerError(message, err))
}

func (s *send) sendResponse(response ApiResponse) {
	s.context.JSON(int(response.GetValue().Status), response)
}

func (s *send) sendError(err error) {
	var res ApiResponse
	var apiError ApiError

	if errors.As(err, &apiError) {
		e := apiError.GetValue()
		switch e.Code {
		case http.StatusBadRequest:
			res = NewBadRequestResponse(e.Message)
		case http.StatusForbidden:
			res = NewForbiddenResponse(e.Message)
		case http.StatusUnauthorized:
			res = NewUnauthorizedResponse(e.Message)
		case http.StatusNotFound:
			res = NewNotFoundResponse(e.Message)
		case http.StatusInternalServerError:
			if s.debug {
				res = NewInternalServerErrorResponse(apiError.Unwrap().Error())
			}
		default:
			if s.debug {
				res = NewInternalServerErrorResponse(apiError.Unwrap().Error())
			}
		}
	}

	if res == nil {
		res = NewInternalServerErrorResponse("An unexpected error occurred. Please try again later.")
	}

	s.sendResponse(res)
}
