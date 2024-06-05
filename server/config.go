package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/unusualcodeorg/go-lang-backend-architecture/internal/utils"
)

var Config config

type config struct {
	// server
	SERVER_HOST string
	SERVER_PORT uint16
	// database
	DB_HOST          string
	DB_PORT          string
	DB_USER          string
	DB_USER_PWD      string
	DB_MIN_POOL_SIZE uint8
	DB_MAX_POOL_SIZE uint8
}

func init() {
	loadEnvFile()
	Config = config{
		// server
		SERVER_HOST: os.Getenv("SERVER_HOST"),
		SERVER_PORT: utils.ParseUint16(os.Getenv("SERVER_PORT")),
		// database
		DB_HOST:          os.Getenv("DB_HOST"),
		DB_PORT:          os.Getenv("DB_PORT"),
		DB_USER:          os.Getenv("DB_USER"),
		DB_USER_PWD:      os.Getenv("DB_USER_PWD"),
		DB_MIN_POOL_SIZE: utils.ParseUint8(os.Getenv("DB_MIN_POOL_SIZE")),
		DB_MAX_POOL_SIZE: utils.ParseUint8(os.Getenv("DB_MAX_POOL_SIZE")),
	}

}

func loadEnvFile() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
