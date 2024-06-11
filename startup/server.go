package startup

import (
	"time"

	"github.com/unusualcodeorg/go-lang-backend-architecture/api/auth"
	authMW "github.com/unusualcodeorg/go-lang-backend-architecture/api/auth/middleware"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/contact"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/user"
	"github.com/unusualcodeorg/go-lang-backend-architecture/config"
	coreMW "github.com/unusualcodeorg/go-lang-backend-architecture/core/middleware"
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
	router.RegisterValidationParsers(network.CustomTagNameFunc())
	dbQueryTimeout := time.Duration(env.DBQueryTimeout) * time.Second

	userService := user.NewUserService(db, dbQueryTimeout)
	authService := auth.NewAuthService(db, dbQueryTimeout, env, userService)
	contactService := contact.NewContactService(db, dbQueryTimeout)

	router.LoadRootMiddlewares(
		coreMW.NewErrorProcessor(), // NOTE: this should be the first handler to be mounted
		authMW.NewKeyProtection(authService),
		coreMW.NewNotFound(),
	)

	authProvider := authMW.NewAuthenticationProvider(authService)
	authorizeProvider := authMW.NewAuthorizationProvider()

	router.LoadControllers(
		auth.NewAuthController(authProvider, authorizeProvider, authService),
		user.NewProfileController(authProvider, authorizeProvider, userService),
		contact.NewContactController(authProvider, authorizeProvider, contactService),
	)

	router.Start(env.ServerHost, env.ServerPort)
}
