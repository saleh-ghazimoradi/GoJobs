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
	JWT          JWT
	UploadDIR    UploadDIR
}

type JWT struct {
	SecretKEY string `env:"JWT_SECRET"`
}

type UploadDIR struct {
	Upload string `env:"UPLOAD_DIR"`
}

type ServerConfig struct {
	Port         string        `env:"SERVER_PORT,required"`
	Version      string        `env:"SERVER_VERSION,required"`
	IdleTimeout  time.Duration `env:"SERVER_IDLE_TIMEOUT,required"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT,required"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT,required"`
}

type DBConfig struct {
	DbHost       string        `env:"DB_HOST,required"`
	DbPort       string        `env:"DB_PORT,required"`
	DbUser       string        `env:"DB_USER,required"`
	DbPassword   string        `env:"DB_PASSWORD,required"`
	DbName       string        `env:"DB_NAME,required"`
	DbSslMode    string        `env:"DB_SSLMODE,required"`
	MaxOpenConns int           `env:"DB_MAX_OPEN_CONNECTIONS,required"`
	MaxIdleConns int           `env:"DB_MAX_IDLE_CONNECTIONS,required"`
	MaxIdleTime  time.Duration `env:"DB_MAX_IDLE_TIME,required"`
	Timeout      time.Duration `env:"DB_TIMEOUT,required"`
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

	jwtConfig := &JWT{}

	if err := env.Parse(jwtConfig); err != nil {
		log.Fatalf("unable to parse config: %v", err)
	}

	config.JWT = *jwtConfig

	uploadDirConfig := &UploadDIR{}
	if err := env.Parse(uploadDirConfig); err != nil {
		log.Fatalf("unable to parse config: %v", err)
	}
	config.UploadDIR = *uploadDirConfig

	AppConfig = config

	return nil
}
