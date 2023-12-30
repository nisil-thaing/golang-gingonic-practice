package main

import (
	"github.com/nisil-thaing/golang-gingonic-practice/internal/api"
	"log"
)

func main() {
	router := api.SetupAPI()
	if err := router.Run(); err != nil {
		log.Fatal("Oops! Couldn't starting the server:", err)
	}
}
