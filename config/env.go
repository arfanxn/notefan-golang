package config

import (
	"github.com/joho/godotenv"
)

func LoadENV() error {
	return godotenv.Load("local.env")
}

func LoadTestENV() error {
	return godotenv.Load("test.env")
}
