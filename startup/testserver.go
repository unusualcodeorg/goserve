package startup

import (
	"net/http/httptest"

	"github.com/unusualcodeorg/goserve/arch/network"
	"github.com/unusualcodeorg/goserve/config"
)

type Teardown = func()

func TestServer() (network.Router, network.Module[module], Teardown) {
	env := config.NewEnv("../.test.env")
	router, module, shutdown := create(env)
	ts := httptest.NewServer(router.GetEngine())
	teardown := func() {
		ts.Close()
		shutdown()
	}
	return router, module, teardown
}
