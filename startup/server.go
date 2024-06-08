package startup

import (
	"time"

	"github.com/unusualcodeorg/go-lang-backend-architecture/api/contact"
	"github.com/unusualcodeorg/go-lang-backend-architecture/config"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/middleware"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

func Server() {
	env := config.NewEnv(".env")

	db := mongo.NewDatabase(env)
	defer db.Disconnect()
	db.Connect()

	router := network.NewRouter(env.GoMode)
	dbQueryTimeout := time.Duration(env.DBQueryTimeout) * time.Second

	router.LoadRootMiddlewares(
		middleware.NewKeyProtection(core.NewContactService(db, dbQueryTimeout)),
		middleware.NewNotFound(),
		middleware.NewErrorHandler(),
	)

	router.LoadControllers(
		contact.NewContactController(
			middleware.NewAuthentication,
			middleware.NewAuthorization,
			contact.NewContactService(db, dbQueryTimeout),
		),
	)

	router.Start(env.ServerHost, env.ServerPort)
}
