package startup

import (
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/contact"
	"github.com/unusualcodeorg/go-lang-backend-architecture/config"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/middleware"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

func Server() {
	env := config.NewEnv(".env")

	db := mongo.NewDatabase(env)
	db.Connect()

	router := network.NewRouter()


	router.LoadControllers(
		contact.NewContactController(
			network.NewBaseController("/contact", middleware.NewAuthenticationMiddleware(), middleware.NewAuthorizationMiddleware()),
			contact.NewService(db),
		),
	)

	router.LoadRootMiddlewares(
		middleware.NewNotFoundMiddleware(),
		middleware.NewApikeyMiddleware(),
	)

	router.Start(env.ServerHost, env.ServerPort)

	defer db.Disconnect()
}
