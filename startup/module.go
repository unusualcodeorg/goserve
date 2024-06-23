package startup

import (
	"context"

	"github.com/unusualcodeorg/goserve/api/auth"
	authMW "github.com/unusualcodeorg/goserve/api/auth/middleware"
	"github.com/unusualcodeorg/goserve/api/blog"
	"github.com/unusualcodeorg/goserve/api/contact"
	"github.com/unusualcodeorg/goserve/api/user"
	coreMW "github.com/unusualcodeorg/goserve/arch/middleware"
	"github.com/unusualcodeorg/goserve/arch/mongo"
	"github.com/unusualcodeorg/goserve/arch/network"
	"github.com/unusualcodeorg/goserve/arch/redis"
	"github.com/unusualcodeorg/goserve/config"
)

type module struct {
	context     context.Context
	env         *config.Env
	db          mongo.Database
	store       redis.Store
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
		blog.NewController(m.AuthenticationProvider(), m.AuthorizationProvider(), blog.NewService(m.db, m.store, m.userService)),
		blog.NewWriterController(m.AuthenticationProvider(), m.AuthorizationProvider(), blog.NewService(m.db, m.store, m.userService)),
		blog.NewEditorController(m.AuthenticationProvider(), m.AuthorizationProvider(), blog.NewService(m.db, m.store, m.userService)),
		contact.NewController(m.AuthenticationProvider(), m.AuthorizationProvider(), contact.NewService(m.db)),
	}
}

func (m *module) RootMiddlewares() []network.RootMiddleware {
	return []network.RootMiddleware{
		coreMW.NewErrorCatcher(), // NOTE: this should be the first handler to be mounted
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

func NewModule(context context.Context, env *config.Env, db mongo.Database, store redis.Store) network.Module[module] {
	userService := user.NewService(db)
	authService := auth.NewService(db, env, userService)

	return &module{
		context:     context,
		env:         env,
		db:          db,
		store:       store,
		userService: userService,
		authService: authService,
	}
}
