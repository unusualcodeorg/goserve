package config

import (
	"github.com/spf13/viper"
	"log"
)

type Env struct {
	GoEnv         string `mapstructure:"GO_ENV"`
	ServerHost    string `mapstructure:"SERVER_HOST"`
	ServerPort    uint16 `mapstructure:"SERVER_PORT"`
	DBHost        string `mapstructure:"DB_HOST"`
	DBName        string `mapstructure:"DB_NAME"`
	DBPort        uint16 `mapstructure:"DB_PORT"`
	DBUser        string `mapstructure:"DB_USER"`
	DBUserPwd     string `mapstructure:"DB_USER_PWD"`
	DBMinPoolSize uint8  `mapstructure:"DB_MIN_POOL_SIZE"`
	DBMaxPoolSize uint8  `mapstructure:"DB_MAX_POOL_SIZE"`
}

func NewEnv(filename string) *Env {
	env := Env{}
	viper.SetConfigFile(filename)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Error reading environment file", err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Error loading environment file", err)
	}

	return &env
}
