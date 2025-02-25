package main

import (
	"fmt"
	"go-markdown-parser/controller"
	"go-markdown-parser/database"
	"go-markdown-parser/routes"
	"go-markdown-parser/utils"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	PORT := os.Getenv("PORT")

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	if PORT == "" {
		PORT = "8080"
	}

	log.Printf("Starting server initialization on port %s", PORT)

	database.StartDB()

	router := gin.New() // Use New() instead of Default() for cleaner logs
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     utils.GetCorsOrigins(),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept", "X-Requested-With"},
		AllowCredentials: true,
	}))

	// Test route
	router.GET("/ping", controller.Ping())
	routes.AuthRoutes(router)
	routes.MarkdownParserRoutes(router)

	serverAddr := fmt.Sprintf("0.0.0.0:%s", PORT)
	log.Printf("Server attempting to listen on %s", serverAddr)

	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
