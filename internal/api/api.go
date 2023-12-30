package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupAPI(uri string) {
	router := gin.Default()

	// Welcome Announcement
	router.GET("/", welcomeController)

	if err := router.Run(uri); err != nil {
		log.Fatal("Oops! Couldn't starting the server:", err)
	}
}

func welcomeController(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to this JWT practice project using Golang and GinGonic!"})
}
