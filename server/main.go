package main

import (
	"fmt"

	"github.com/unusualcodeorg/go-lang-backend-architecture/api"
	"github.com/unusualcodeorg/go-lang-backend-architecture/config"
	"github.com/unusualcodeorg/go-lang-backend-architecture/internal/core"
)

func main() {
	defer core.DisconnectMongoDb()
	address := fmt.Sprintf("%s:%d", config.Env.SERVER_HOST, config.Env.SERVER_PORT)
	api.CreateRouter().Run(address)
}
