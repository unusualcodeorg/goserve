package main

import "github.com/unusualcodeorg/go-lang-backend-architecture/api"

func main() {
	api.CreateRouter().Run("localhost:8080")
}
