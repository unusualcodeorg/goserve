package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReqBody(t *testing.T) {
	body := `{"field": "test"}`
	ctx, _ := MockTestHttp(t, "POST", "/mock", "/mock", body, MockSuccessMsgHandler("success"))

	dto, err := ReqBody(ctx, &MockDto{})

	assert.NoError(t, err)
	assert.Equal(t, dto.Field, "test")
}

func TestReqBody_Error(t *testing.T) {
	body := `{"wrong": "test"}`
	ctx, _ := MockTestHttp(t, "POST", "/mock", "/mock", body, MockSuccessMsgHandler("success"))

	dto, err := ReqBody(ctx, &MockDto{})

	assert.Nil(t, dto)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "field is required")
}

func TestReqQuery(t *testing.T) {
	ctx, _ := MockTestHttp(t, "GET", "/mock", "/mock?field=test", "", MockSuccessMsgHandler("success"))

	dto, err := ReqQuery(ctx, &MockDto{})

	assert.NoError(t, err)
	assert.Equal(t, dto.Field, "test")
}

func TestReqQuery_Error(t *testing.T) {
	ctx, _ := MockTestHttp(t, "GET", "/mock", "/mock?wrong=test", "", MockSuccessMsgHandler("success"))

	dto, err := ReqQuery(ctx, &MockDto{})

	assert.Nil(t, dto)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "field is required")
}
