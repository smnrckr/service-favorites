package utils

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnviromentVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}
