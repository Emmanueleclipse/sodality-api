package middlewares

import (
	"log"
	"os"
	"path"

	"github.com/joho/godotenv"
)

// DotEnvVariable -> get .env
func DotEnvVariable(key string) string {

	// load .env file
	wd, _ := os.Getwd()
	log.Println(path.Join(wd, "/.env"))

	err := godotenv.Load()

	if err != nil {
		log.Printf("Error loading .env file %s", err)
	}

	return os.Getenv(key)
}
