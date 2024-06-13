package startup

import (
	"context"

	"github.com/unusualcodeorg/go-lang-backend-architecture/api/auth"
	authMW "github.com/unusualcodeorg/go-lang-backend-architecture/api/auth/middleware"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/blog"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/contact"
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/user"
	"github.com/unusualcodeorg/go-lang-backend-architecture/config"
	coreMW "github.com/unusualcodeorg/go-lang-backend-architecture/core/middleware"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

type module struct {
	context     context.Context
	env         *config.Env
	db          mongo.Database
	userService user.Service
	authService auth.Service
}

func (m *module) GetInstance() *module {
	return m
}

func (m *module) Controllers() []network.Controller {
	return []network.Controller{
		auth.NewController(m.AuthenticationProvider(), m.AuthorizationProvider(), m.authService),
		user.NewController(m.AuthenticationProvider(), m.AuthorizationProvider(), m.userService),
		blog.NewController(m.AuthenticationProvider(), m.AuthorizationProvider(), blog.NewService(m.db)),
		blog.NewWriterController(m.AuthenticationProvider(), m.AuthorizationProvider(), blog.NewService(m.db)),
		blog.NewWriterController(m.AuthenticationProvider(), m.AuthorizationProvider(), blog.NewService(m.db)),
		contact.NewController(m.AuthenticationProvider(), m.AuthorizationProvider(), contact.NewService(m.db)),
	}
}

func (m *module) RootMiddlewares() []network.RootMiddleware {
	return []network.RootMiddleware{
		coreMW.NewErrorProcessor(), // NOTE: this should be the first handler to be mounted
		authMW.NewKeyProtection(m.authService),
		coreMW.NewNotFound(),
	}
}

func (m *module) AuthenticationProvider() network.AuthenticationProvider {
	return authMW.NewAuthenticationProvider(m.authService, m.userService)
}

func (m *module) AuthorizationProvider() network.AuthorizationProvider {
	return authMW.NewAuthorizationProvider()
}

func NewModule(context context.Context, env *config.Env, db mongo.Database) network.Module[module] {
	userService := user.NewService(db)
	authService := auth.NewService(db, env, userService)

	return &module{
		context:     context,
		env:         env,
		db:          db,
		userService: userService,
		authService: authService,
	}
}
