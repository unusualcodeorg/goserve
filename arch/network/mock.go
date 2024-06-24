package network

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type MockDto struct {
	Field string `json:"field" form:"field" uri:"field" binding:"required" validate:"required"`
}

func (d *MockDto) GetValue() *MockDto {
	return d
}

func (b *MockDto) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	var msgs []string
	for _, err := range errs {
		switch err.Tag() {
		case "required":
			msgs = append(msgs, fmt.Sprintf("%s is required", err.Field()))
		default:
			msgs = append(msgs, fmt.Sprintf("%s is invalid", err.Field()))
		}
	}
	return msgs, nil
}

func MockSuccessMsgHandler(msg string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		NewResponseSender().Send(ctx).SuccessMsgResponse(msg)
	}
}

func MockSuccessDataHandler(msg string, data any) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		NewResponseSender().Send(ctx).SuccessDataResponse(msg, data)
	}
}

func MockTestHttp(
	t *testing.T, httpMethod, path, url, body string,
	handler gin.HandlerFunc,
) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)

	rr := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(rr)
	r.Handle(httpMethod, path, handler)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(CustomTagNameFunc())
	}

	req, err := http.NewRequest(httpMethod, url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	ctx.Request = req

	r.ServeHTTP(rr, req)

	return ctx, rr
}

func MockTestRootMiddleware(
	t *testing.T, httpMethod, path, url, body string,
	m RootMiddleware,
	handler gin.HandlerFunc,
) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)

	rr := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(rr)
	m.Attach(r)
	r.Handle(httpMethod, path, handler)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(CustomTagNameFunc())
	}

	req, err := http.NewRequest(httpMethod, url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	ctx.Request = req

	r.ServeHTTP(rr, req)

	return ctx, rr
}
