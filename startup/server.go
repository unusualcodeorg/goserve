package startup

import (
	"github.com/unusualcodeorg/go-lang-backend-architecture/api/contact"
	"github.com/unusualcodeorg/go-lang-backend-architecture/common/mongo"
	"github.com/unusualcodeorg/go-lang-backend-architecture/common/network"
	"github.com/unusualcodeorg/go-lang-backend-architecture/config"
)

func Server() {
	env := config.NewEnv(".env")

	db := mongo.NewDatabase(env)
	db.Connect()

	contactService := contact.NewService(db)
	contactController := contact.NewContactController(contactService)

	router := network.NewRouter()

	router.LoadControllers(contactController)
	router.LoadHandlers(network.NotFoundHandler)

	router.Start(env.ServerHost, env.ServerPort)

	defer db.Disconnect()
}
