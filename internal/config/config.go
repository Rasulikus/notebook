package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	ErrSecretMissing = errors.New("SECRET is missing or empty")
)

const (
	keyDBHost, defaultDBHost = "DB_HOST", "localhost"
	keyDBPort, defaultDBPort = "DB_PORT", "5432"
	keyDBUser, defaultDBUser = "DB_USER", "root"
	keyDBPass, defaultDBPass = "DB_PASS", "example"
	keyDBName, defaultDBName = "DB_NAME", "notebook"

	keyAuthSecret, defaultAuthSecret = "AUTH_SECRET", "1420061f070b81aac84ceb449812770ab9d1f1d6b4c0aba33533ce6dde6f96fb"

	LogDefaultValue = "%s is missing, using default value"
)

type Config struct {
	Db   DbConfig
	Auth AuthConfig
}

type DbConfig struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

type AuthConfig struct {
	Secret string
}

func getEnv(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		log.Printf(LogDefaultValue, key)
		return defaultValue
	}
	return value
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("load .env: %v", err)
	}

	cfg := new(Config)

	cfg.Db.Host = getEnv(keyDBHost, defaultDBHost)
	cfg.Db.Port = getEnv(keyDBPort, defaultDBPort)
	cfg.Db.User = getEnv(keyDBUser, defaultDBUser)
	cfg.Db.Pass = getEnv(keyDBPass, defaultDBPass)
	cfg.Db.Name = getEnv(keyDBName, defaultDBName)

	cfg.Auth.Secret = getEnv(keyAuthSecret, defaultAuthSecret)

	return cfg
}
