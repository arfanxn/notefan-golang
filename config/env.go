package config

import (
	"github.com/joho/godotenv"
)

func InitializeENV() error {
	return godotenv.Load("local.env")
}
