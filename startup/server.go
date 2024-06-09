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

	EnsureDbIndexes(db)

	router := network.NewRouter(env.GoMode)
	dbQueryTimeout := time.Duration(env.DBQueryTimeout) * time.Second

	router.LoadRootMiddlewares(
		middleware.NewErrorHandler(), // NOTE: this should be the first handler to be mounted
		middleware.NewKeyProtection(core.NewContactService(db, dbQueryTimeout)),
		middleware.NewNotFound(),
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
