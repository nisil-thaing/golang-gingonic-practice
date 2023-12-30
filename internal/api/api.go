package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupAPI() *gin.Engine {
	router := gin.Default()

	// Welcome Announcement
	router.GET("/", welcomeController)

	return router
}

func welcomeController(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to this JWT practice project using Golang and GinGonic!"})
}
