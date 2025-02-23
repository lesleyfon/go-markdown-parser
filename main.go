package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"main.go/api"
)

func getCorsOrigins() []string {
	// Default origins for development
	var CORS_ALLOWED_ORIGINS []string

	// Try to get origins from env
	ALLOW_ORIGINS_JSON := os.Getenv("ALLOW_ORIGINS")
	if ALLOW_ORIGINS_JSON == "" {
		log.Println("Warning: ALLOW_ORIGINS not set, using default origins")
		return CORS_ALLOWED_ORIGINS
	}

	if err := json.Unmarshal([]byte(ALLOW_ORIGINS_JSON), &CORS_ALLOWED_ORIGINS); err != nil {
		log.Printf("Error parsing CORS origins: %v. Using default origins", err)
		return CORS_ALLOWED_ORIGINS
	}

	if len(CORS_ALLOWED_ORIGINS) == 0 {
		log.Println("Warning: No CORS origins provided, using default origins")
		return CORS_ALLOWED_ORIGINS
	}

	return CORS_ALLOWED_ORIGINS
}

func main() {
	PORT := os.Getenv("PORT")

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}
	if PORT == "" {
		PORT = "10000"
	}

	// Configure CORS
	corsConfig := cors.Config{
		AllowOrigins:     getCorsOrigins(),
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	// Use the shared router setup with our custom CORS config
	router := api.SetupRouter(&corsConfig)

	// Add logger middleware for local development
	router.Use(gin.Logger())

	log.Printf("Server starting on port %s", PORT)
	router.Run(fmt.Sprintf("0.0.0.0:%s", PORT))
}
