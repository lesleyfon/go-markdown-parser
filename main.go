package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"main.go/database"
	"main.go/routes"
)

func main() {
	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = "8080"
	}

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}
	database.StartDB()

	router := gin.Default()
	router.Use(gin.Logger())
	// Register Routes
	routes.AuthRoutes(router)
	routes.MarkdownParserRoutes(router)
	// Test routes
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK,
			gin.H{
				"message": "Welcome to Markdown parser",
			})
	})

	router.Run()
}
