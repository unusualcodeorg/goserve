package mongo

import (
	"context"
	"fmt"
	"log"

	"github.com/unusualcodeorg/go-lang-backend-architecture/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database interface {
	Client() *mongo.Client
	GetCollection(string) *mongo.Collection
	Connect()
	Disconnect()
}

type database struct {
	user        string
	pwd         string
	host        string
	port        uint16
	name        string
	minPoolSize uint16
	maxPoolSize uint16
	client      *mongo.Client
}

func NewDatabase(env *config.Env) Database {
	db := database{
		user:        env.DBUser,
		pwd:         env.DBUserPwd,
		host:        env.DBHost,
		port:        env.DBPort,
		name:        env.DBName,
		minPoolSize: env.DBMinPoolSize,
		maxPoolSize: env.DBMaxPoolSize,
	}
	return &db
}

func (db *database) Connect() {
	ctx := context.TODO()

	uri := fmt.Sprintf(
		"mongodb://%s:%s@%s:%d/%s",
		db.user, db.pwd, db.host, db.port, db.name,
	)

	clientOptions := options.Client().ApplyURI(uri)
	clientOptions.SetMaxPoolSize(uint64(db.maxPoolSize))
	clientOptions.SetMaxPoolSize(uint64(db.minPoolSize))

	fmt.Println("Connecting Mongo...")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Connection to Mongo Failed!: ", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Panic("Pinging to Mongo Failed!: ", err)
	}
	fmt.Println("Connected to Mongo!")

	db.client = client
}

func (db *database) GetCollection(name string) *mongo.Collection {
	return db.client.Database(db.name).Collection(name)
}

func (db *database) Disconnect() {
	fmt.Println("Disconnecting Mongo...")
	err := db.client.Disconnect(context.TODO())
	if err != nil {
		log.Panic(err)
	}
}

func (db *database) Client() *mongo.Client {
	return db.client
}
