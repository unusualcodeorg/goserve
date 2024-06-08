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

	dbQuery := mongo.NewDatabaseQuery(db)

	router := network.NewRouter(env.GoMode)

	router.LoadControllers(
		contact.NewContactController(middleware.NewAuthentication, middleware.NewAuthorization, contact.NewContactService(dbQuery)),
	)

	router.LoadRootMiddlewares(
		middleware.NewKeyProtection(),
		middleware.NewNotFound(),
	)

	router.Start(env.ServerHost, env.ServerPort)

	defer db.Disconnect()
}
