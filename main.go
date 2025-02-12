package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
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
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
