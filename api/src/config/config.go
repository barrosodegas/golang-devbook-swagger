package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	ConectionDataBaseString = ""
	Port                    = 0
	SecretKey               []byte
)

// Load loads the system environment variables.
func Load() {
	var error error
	var envFilePath = os.Getenv("ENV_PATH")

	if envFilePath == "" {
		envFilePath = ".env"
	}

	if error = godotenv.Load(envFilePath); error != nil {
		log.Fatal(error)
	}

	Port, error = strconv.Atoi(os.Getenv("API_PORT"))

	if error != nil {
		Port = 9000
	}

	ConectionDataBaseString = fmt.Sprintf(
		"%s:%s@/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	SecretKey = []byte(os.Getenv("SECRET_KEY"))
}
