// Code coverage is disabled for main package as it's the application entry point
//go:build skip_coverage

// @title Internal CS Center API
// @version 1.0
// @description This is the API server for Internal CS Center
// @host localhost:8080
// @BasePath /api/v1

package main

import (
	"log"

	// IMPORTANT: import docs as a named package so we can override fields at runtime
	_ "abasithdev.github.io/internal-cs-center-backend/docs"
	"abasithdev.github.io/internal-cs-center-backend/internal/config"
	"abasithdev.github.io/internal-cs-center-backend/internal/router"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load() // loads .env if present

	config.ConfigureSwagger()
	appConfig := config.Load()

	log.Println("Starting server :" + appConfig.Port)

	r := router.NewRouter()
	if err := r.Run(":" + appConfig.Port); err != nil {
		log.Fatal(err)
	}
}
