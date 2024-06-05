package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var Config config

type config struct {
	// server
	SERVER_HOST string
	SERVER_PORT string
	// database
	DB_HOST     string
	DB_PORT     string
	DB_USER     string
	DB_PASSWORD string
}

func init() {
	loadEnvFile()
	Config = config{
		// server
		SERVER_HOST: os.Getenv("SERVER_HOST"),
		SERVER_PORT: os.Getenv("SERVER_PORT"),
		// database
		DB_HOST:     os.Getenv("DB_HOST"),
		DB_PORT:     os.Getenv("DB_PORT"),
		DB_USER:     os.Getenv("DB_USER"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
	}
}

func loadEnvFile() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
