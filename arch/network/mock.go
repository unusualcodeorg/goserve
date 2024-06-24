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
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func MockTestHandler(
	t *testing.T, httpMethod, path, url, body string,
	handler gin.HandlerFunc,
	headers ...primitive.E,
) *httptest.ResponseRecorder {
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

	for _, h := range headers {
		req.Header.Set(h.Key, h.Value.(string))
	}

	ctx.Request = req

	r.ServeHTTP(rr, req)

	return rr
}

func MockTestRootMiddleware(
	t *testing.T, httpMethod, path, url, body string,
	middleware RootMiddleware,
	handler gin.HandlerFunc,
	headers ...primitive.E,
) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)

	rr := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(rr)
	middleware.Attach(r)
	r.Handle(httpMethod, path, handler)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(CustomTagNameFunc())
	}

	req, err := http.NewRequest(httpMethod, url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	for _, h := range headers {
		req.Header.Set(h.Key, h.Value.(string))
	}

	ctx.Request = req

	r.ServeHTTP(rr, req)

	return rr
}

func MockTestAuthenticationProvider(
	t *testing.T, httpMethod, path, url, body string,
	auth AuthenticationProvider,
	handler gin.HandlerFunc,
	headers ...primitive.E,
) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)

	rr := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(rr)
	r.Use(auth.Middleware())
	r.Handle(httpMethod, path, handler)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(CustomTagNameFunc())
	}

	req, err := http.NewRequest(httpMethod, url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	for _, h := range headers {
		req.Header.Set(h.Key, h.Value.(string))
	}

	ctx.Request = req

	r.ServeHTTP(rr, req)

	return rr
}

func MockTestAuthorizationProvider(
	t *testing.T, httpMethod, path, url, body string,
	auth AuthenticationProvider,
	authz AuthenticationProvider,
	handler gin.HandlerFunc,
	headers ...primitive.E,
) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)

	rr := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(rr)
	r.Use(auth.Middleware())
	r.Use(authz.Middleware())
	r.Handle(httpMethod, path, handler)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(CustomTagNameFunc())
	}

	req, err := http.NewRequest(httpMethod, url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	for _, h := range headers {
		req.Header.Set(h.Key, h.Value.(string))
	}

	ctx.Request = req

	r.ServeHTTP(rr, req)

	return rr
}
