package main

import (
	"ibuYemekBotu/services"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting Bot")
	services.TelegramHandler()
	log.Println(os.Getenv("MONGODB_URI"))
}

func init() {
	log.Println(os.Getenv("MONGODB_URI"))
	err := godotenv.Load(".envDev")
	if err != nil {
		log.Println("Error loading .env file")
	}
}
