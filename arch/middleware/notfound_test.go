package middleware

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unusualcodeorg/goserve/arch/network"
)

func TestNotFoundMiddleware(t *testing.T) {
	_, rr := network.MockTestRootMiddleware(t, "GET", "/test", "/no", "",
		NewNotFound(),
		network.MockSuccessMsgHandler("success"),
	)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Contains(t, rr.Body.String(), `"message":"url not found"`)
}
