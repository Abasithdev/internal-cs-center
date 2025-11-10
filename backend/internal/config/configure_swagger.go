package config

import (
	"os"

	"abasithdev.github.io/internal-cs-center-backend/docs"
)

// ConfigureSwagger sets up Swagger configuration from environment variables
func ConfigureSwagger() {
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
}
