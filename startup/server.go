package startup

import (
	"time"

	"github.com/unusualcodeorg/go-lang-backend-architecture/api/auth"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/contact"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/profile"
	"github.com/unusualcodeorg/go-lang-backend-architecture/config"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core"
	m "github.com/unusualcodeorg/go-lang-backend-architecture/core/middleware"
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

	secretService := core.NewSecretService(db, dbQueryTimeout)
	tokenService := core.NewTokenService(db, dbQueryTimeout, env)
	userService := core.NewUserService(db, dbQueryTimeout)
	profileService := profile.NewProfileService(db, dbQueryTimeout)
	authService := auth.NewAuthService(db, dbQueryTimeout, userService, tokenService)
	contactService := contact.NewContactService(db, dbQueryTimeout)

	router.LoadRootMiddlewares(
		m.NewErrorProcessor(), // NOTE: this should be the first handler to be mounted
		m.NewKeyProtection(secretService),
		m.NewNotFound(),
	)

	authProvider := m.NewAuthProvider(userService, tokenService)
	authorizeProvider := m.NewAuthorizeProvider()

	router.LoadControllers(
		auth.NewAuthController(authProvider, authorizeProvider, authService),
		profile.NewProfileController(authProvider, authorizeProvider, userService, profileService),
		contact.NewContactController(authProvider, authorizeProvider, contactService),
	)

	router.Start(env.ServerHost, env.ServerPort)
}
