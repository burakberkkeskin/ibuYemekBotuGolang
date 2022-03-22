package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {

	telegramHandler()

}

func init() {
	err := godotenv.Load(".envDev")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
