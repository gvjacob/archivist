package utils

import (
	"github.com/charmbracelet/log"
	"github.com/lpernett/godotenv"
)

func LoadDotEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Info("Loaded .env file")
}
