package utils

import (
	b64 "encoding/base64"
	"os"

	"github.com/charmbracelet/log"
	"github.com/lpernett/godotenv"
)

func LoadDotEnv() {
	if os.Getenv("RAILWAY_ENVIRONMENT_NAME") == "production" {
		log.Info("Running in production mode")
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Info("Loaded .env file")
}

func GetBasicAuthorizationHeader() string {
	bytes := []byte(os.Getenv("SPOTIFY_CLIENT_ID") + ":" + os.Getenv("SPOTIFY_CLIENT_SECRET"))
	return "Basic " + b64.URLEncoding.EncodeToString(bytes)
}
