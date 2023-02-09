package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

func InitializeENV() error {
	return LoadENV("local.env")
}

func LoadENV(filename ...string) error {
	return godotenv.Load(filename...)
}

func GetENVKey(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", errors.New("Environment key not found: " + key)
	}
	return value, nil
}
