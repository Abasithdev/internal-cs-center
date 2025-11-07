// @title Internal CS Center API
// @version 1.0
// @description Internal API for CS center.
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Type "Bearer <token>" to authenticate.

// @host localhost:8080
// @BasePath /dashboard/v1
package main

import (
	"log"
	"os"

	// IMPORTANT: import docs as a named package so we can override fields at runtime
	"abasithdev.github.io/internal-cs-center-backend/cmd/server/docs"
	"abasithdev.github.io/internal-cs-center-backend/internal/router"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load() // loads .env if present

	// runtime overrides (fallback to generated values if env not set)
	if h := os.Getenv("SWAGGER_HOST"); h != "" {
		docs.SwaggerInfo.Host = h
	}
	if b := os.Getenv("SWAGGER_BASEPATH"); b != "" {
		docs.SwaggerInfo.BasePath = b
	}
	if t := os.Getenv("SWAGGER_TITLE"); t != "" {
		docs.SwaggerInfo.Title = t
	}
	if v := os.Getenv("SWAGGER_VERSION"); v != "" {
		docs.SwaggerInfo.Version = v
	}
	if d := os.Getenv("SWAGGER_DESCRIPTION"); d != "" {
		docs.SwaggerInfo.Description = d
	}

	log.Println("Starting server :8080")

	r := router.NewRouter()
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
