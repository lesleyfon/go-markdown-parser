package utils

import (
	"encoding/json"
	"log"
	"os"
)

func GetCorsOrigins() []string {
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

	log.Printf("CORS origins: %v", CORS_ALLOWED_ORIGINS)

	return CORS_ALLOWED_ORIGINS
}
