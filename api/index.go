package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"main.go/database"
	"main.go/routes" // Update this to your module name
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

// Handler - Vercel entry point
func Handler(w http.ResponseWriter, r *http.Request) {
	// Set Gin to release mode in production
	gin.SetMode(gin.ReleaseMode)

	// Initialize router
	router := SetupRouter(nil) // Pass nil for default CORS config

	// Serve the request
	router.ServeHTTP(w, r)
}

// SetupRouter initializes and returns the Gin router
func SetupRouter(corsConfig *cors.Config) *gin.Engine {
	// Initialize database
	database.StartDB()

	// Create router
	router := gin.New() // Use New() instead of Default() for serverless

	// Configure CORS
	if corsConfig == nil {
		defaultConfig := cors.DefaultConfig()
		defaultConfig.AllowOrigins = getCorsOrigins()
		defaultConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
		defaultConfig.AllowHeaders = []string{"Origin", "Authorization", "Content-Type"}
		defaultConfig.AllowCredentials = true
		corsConfig = &defaultConfig
	}
	router.Use(cors.New(*corsConfig))

	// Setup routes
	routes.AuthRoutes(router)
	routes.MarkdownParserRoutes(router)

	// Add a health check route
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "API is running",
		})
	})

	return router
}
