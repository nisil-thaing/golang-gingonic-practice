package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nisil-thaing/golang-gingonic-practice/internal/app/auth"
)

func SetupAPI(uri string) {
	router := gin.Default()

	publicRouter := router.Group("/api")
	// Welcome Announcement
	publicRouter.GET("/", welcomeController)

	auth.SetupRoutes(publicRouter.Group("/auth"))

	if err := router.Run(uri); err != nil {
		log.Fatal("Oops! Couldn't starting the server:", err)
	}
}

func welcomeController(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to this JWT practice project using Golang and GinGonic!"})
}
