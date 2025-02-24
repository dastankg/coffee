package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Db   DbConfig
	Auth AuthConfig
}

type DbConfig struct {
	DATABASE_URL string
}

type AuthConfig struct {
	AccessSecret  string
	RefreshSecret string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file, using default config")
	}

	return &Config{
		Db: DbConfig{
			DATABASE_URL: os.Getenv("DATABASE_URL"),
		},
		Auth: AuthConfig{
			AccessSecret:  os.Getenv("TOKEN"),
			RefreshSecret: os.Getenv("TOKEN"),
		},
	}
}
