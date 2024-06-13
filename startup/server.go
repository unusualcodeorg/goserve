package startup

import (
	"context"
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
	ctx := context.Background()
	env := config.NewEnv(".env")

	dbConfig := mongo.DbConfig{
		User:        env.DBUser,
		Pwd:         env.DBUserPwd,
		Host:        env.DBHost,
		Port:        env.DBPort,
		Name:        env.DBName,
		MinPoolSize: env.DBMinPoolSize,
		MaxPoolSize: env.DBMaxPoolSize,
		Timeout: time.Duration(env.DBQueryTimeout) * time.Second,
	}

	db := mongo.NewDatabase(ctx, dbConfig)
	defer db.Disconnect()
	db.Connect()

	EnsureDbIndexes(db)

	router := network.NewRouter(env.GoMode)
	router.RegisterValidationParsers(network.CustomTagNameFunc())

	userService := user.NewUserService(db)
	authService := auth.NewAuthService(db, env, userService)
	contactService := contact.NewContactService(db)

	router.LoadRootMiddlewares(
		coreMW.NewErrorProcessor(), // NOTE: this should be the first handler to be mounted
		authMW.NewKeyProtection(authService),
		coreMW.NewNotFound(),
	)

	authProvider := authMW.NewAuthenticationProvider(authService, userService)
	authorizeProvider := authMW.NewAuthorizationProvider()

	router.LoadControllers(
		auth.NewAuthController(authProvider, authorizeProvider, authService),
		user.NewProfileController(authProvider, authorizeProvider, userService),
		contact.NewContactController(authProvider, authorizeProvider, contactService),
	)

	router.Start(env.ServerHost, env.ServerPort)
}
