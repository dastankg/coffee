package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Db   DbConfig
	Auth AuthConfig
	Smtp SmtpConfig
}

type SmtpConfig struct {
	SmtpHost string
	SmtpPort string
	From     string
	Password string
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
		Smtp: SmtpConfig{
			SmtpHost: os.Getenv("SMTP_HOST"),
			SmtpPort: os.Getenv("SMTP_PORT"),
			From:     os.Getenv("SMTP_EMAIL"),
			Password: os.Getenv("SMTP_PASSWORD"),
		},
	}
}
