package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	UrlApi  = ""
	PortApp = 0
	// Used to authenticate cookie
	HashKey []byte
	// Used to encrypt cookie data
	BlockKey []byte
)

func LoadVars() {
	var error error

	if error = godotenv.Load(); error != nil {
		log.Fatal(error)
	}

	PortApp, error = strconv.Atoi(os.Getenv("PORT_APP"))
	if error != nil {
		log.Fatal(error)
	}

	UrlApi = os.Getenv("URL_API")
	HashKey = []byte(os.Getenv("HASH_KEY"))
	BlockKey = []byte(os.Getenv("BLOCK_KEY"))
}
