package core

import (
	"context"
	"fmt"
	"log"

	"github.com/unusualcodeorg/go-lang-backend-architecture/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client mongo.Client

func MongoDatabase() *mongo.Database {
	return client.Database(config.Env.DB_NAME)
}

func init() {
	// Set client options
	uri := fmt.Sprintf(
		"mongodb://%s:%s@%s:%d/%s",
		config.Env.DB_USER, config.Env.DB_USER_PWD, config.Env.DB_HOST, config.Env.DB_PORT, config.Env.DB_NAME,
	)

	fmt.Println("Connecting to MongoDB:", uri)

	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
		fmt.Println("Failed to connect MongoDB:", err)
	}

	fmt.Println("Connected to MongoDB!")
}

func ConnectMongoDb() {}

func DisconnectMongoDb() {
	fmt.Println("Disconnect to MongoDB!")
	client.Disconnect(context.TODO())
}
