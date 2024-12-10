package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"log"
	"time"
)

var AppConfig *Config

type Config struct {
	ServerConfig ServerConfig
	DBConfig     DBConfig
}

type ServerConfig struct {
	Port         string        `env:"SERVER_PORT,required"`
	Version      string        `env:"SERVER_VERSION,required"`
	IdleTimeout  time.Duration `env:"SERVER_IDLE_TIMEOUT,required"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT,required"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT,required"`
}

type DBConfig struct {
}

func LoadingConfig() error {
	if err := godotenv.Load("app.env"); err != nil {
		log.Fatalf("Unable to load app.env file: %v", err)
	}

	config := &Config{}

	if err := env.Parse(config); err != nil {
		log.Fatalf("unable to parse config: %v", err)
	}

	serverConfig := &ServerConfig{}

	if err := env.Parse(serverConfig); err != nil {
		log.Fatalf("unable to parse config: %v", err)
	}

	config.ServerConfig = *serverConfig

	dbConfig := &DBConfig{}
	if err := env.Parse(dbConfig); err != nil {
		log.Fatalf("unable to parse config: %v", err)
	}

	config.DBConfig = *dbConfig

	AppConfig = config

	return nil
}
