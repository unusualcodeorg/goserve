package startup

import (
	"time"

	"github.com/unusualcodeorg/go-lang-backend-architecture/api/auth"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/contact"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/user"
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

	coreService := core.NewCoreService(db, dbQueryTimeout)
	userService := user.NewUserService(db, dbQueryTimeout)
	authService := auth.NewAuthService(db, dbQueryTimeout, env, userService)
	contactService := contact.NewContactService(db, dbQueryTimeout)

	router.LoadRootMiddlewares(
		middleware.NewErrorHandler(), // NOTE: this should be the first handler to be mounted
		middleware.NewKeyProtection(coreService),
		middleware.NewNotFound(),
	)

	router.LoadControllers(
		auth.NewAuthController(middleware.NewAuthentication, middleware.NewAuthorization, authService),
		user.NewUserController(middleware.NewAuthentication, middleware.NewAuthorization, userService),
		contact.NewContactController(middleware.NewAuthentication, middleware.NewAuthorization, contactService),
	)

	router.Start(env.ServerHost, env.ServerPort)
}
