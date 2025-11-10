package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	JwtSecret      string
	AllowedOrigins []string
}

func Load() *Config {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	rawOrigins := os.Getenv("ALLOWED_ORIGINS")
	var origins []string
	if rawOrigins != "" {
		for _, o := range strings.Split(rawOrigins, ",") {
			o = strings.TrimSpace(o)
			if o != "" {
				origins = append(origins, o)
			}
		}
	} else {
		origins = []string{"*"} // fallback
		log.Println("⚠️  ALLOWED_ORIGINS not set, allowing all (*)")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("⚠️  JWT_SECRET not set, using default")
		secret = "changeme"
	}

	return &Config{
		Port:           port,
		JwtSecret:      secret,
		AllowedOrigins: origins,
	}
}
