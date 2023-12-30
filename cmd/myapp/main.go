package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nisil-thaing/golang-gingonic-practice/internal/api"
)

var DEFAULT_PORT string = "8080"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Oops! Couldn't load the .env file! Please refer to the .env.sample file!")
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = DEFAULT_PORT
	}

	uri := ":" + port
	api.SetupAPI(uri)
}
