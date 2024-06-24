package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReqBody(t *testing.T) {
	body := `{"field": "test"}`
	ctx, _ := MockHttpWithBody(t, "POST", "/mock", "/mock", MockSuccessMsgHandler("success"), body)

	dto, err := ReqBody(ctx, &MockDto{})

	assert.NoError(t, err)
	assert.Equal(t, dto.Field, "test")
}

func TestReqBody_Error(t *testing.T) {
	body := `{"wrong": "test"}`
	ctx, _ := MockHttpWithBody(t, "POST", "/mock", "/mock", MockSuccessMsgHandler("success"), body)

	dto, err := ReqBody(ctx, &MockDto{})

	assert.Nil(t, dto)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "field is required")
}

func TestReqQuery(t *testing.T) {
	ctx, _ := MockHttp(t, "GET", "/mock", "/mock?field=test", MockSuccessMsgHandler("success"))

	dto, err := ReqQuery(ctx, &MockDto{})

	assert.NoError(t, err)
	assert.Equal(t, dto.Field, "test")
}

func TestReqQuery_Error(t *testing.T) {
	ctx, _ := MockHttp(t, "GET", "/mock", "/mock?wrong=test", MockSuccessMsgHandler("success"))

	dto, err := ReqQuery(ctx, &MockDto{})

	assert.Nil(t, dto)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "field is required")
}

