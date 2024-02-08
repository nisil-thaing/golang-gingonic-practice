package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nisil-thaing/golang-gingonic-practice/internal/api"
	"github.com/nisil-thaing/golang-gingonic-practice/internal/database"
)

var DEFAULT_PORT string = "8080"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Oops! Couldn't load the .env file! Please refer to the .env.sample file!")
	}

	port := os.Getenv("PORT")
	mongoDBUri := os.Getenv("MONGODB_URI")
	database.SetupDBConnection(mongoDBUri)

	if port == "" {
		port = DEFAULT_PORT
	}

	uri := ":" + port
	api.SetupAPI(uri)
}
