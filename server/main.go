package main

import (
	"fmt"

	"github.com/unusualcodeorg/go-lang-backend-architecture/api"
)

func main() {
	address := fmt.Sprintf("%s:%s", Config.SERVER_HOST, Config.SERVER_PORT)
	api.CreateRouter().Run(address)
}
