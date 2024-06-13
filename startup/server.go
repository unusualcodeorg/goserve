package startup

import (
	"context"
	"time"

	"github.com/unusualcodeorg/go-lang-backend-architecture/config"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/core/network"
)

func Server() {
	context := context.Background()
	env := config.NewEnv(".env")

	dbConfig := mongo.DbConfig{
		User:        env.DBUser,
		Pwd:         env.DBUserPwd,
		Host:        env.DBHost,
		Port:        env.DBPort,
		Name:        env.DBName,
		MinPoolSize: env.DBMinPoolSize,
		MaxPoolSize: env.DBMaxPoolSize,
		Timeout:     time.Duration(env.DBQueryTimeout) * time.Second,
	}

	db := mongo.NewDatabase(context, dbConfig)
	db.Connect()
	defer db.Disconnect()
	EnsureDbIndexes(db)

	module := NewModule(context, env, db)

	router := network.NewRouter(env.GoMode)
	router.RegisterValidationParsers(network.CustomTagNameFunc())
	router.LoadRootMiddlewares(module.RootMiddlewares())
	router.LoadControllers(module.Controllers())
	router.Start(env.ServerHost, env.ServerPort)
}
