// Code coverage is disabled for router package as it's a thin wrapper around gin
//go:build skip_coverage

package router

import (
	"log"
	"strings"
	"time"

	"abasithdev.github.io/internal-cs-center-backend/internal/config"
	"abasithdev.github.io/internal-cs-center-backend/internal/handler"
	"abasithdev.github.io/internal-cs-center-backend/internal/middleware"
	"abasithdev.github.io/internal-cs-center-backend/internal/service"
	"abasithdev.github.io/internal-cs-center-backend/internal/storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter() *gin.Engine {
	store := storage.NewMemoryStore()
	authService := service.NewAuthService(store, []byte("donttellanyone"))
	paymentService := service.NewPaymentService(store)

	authHandler := handler.NewAuthHandler(authService)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	appConfig := config.Load()

	r := gin.Default()

	// Normalize and validate allowed origins to avoid panics from the CORS middleware
	var allowOrigins []string
	for _, o := range appConfig.AllowedOrigins {
		o = strings.TrimSpace(o)
		if o == "" {
			continue
		}
		if o == "*" || strings.HasPrefix(o, "http://") || strings.HasPrefix(o, "https://") {
			allowOrigins = append(allowOrigins, o)
			continue
		}
		if strings.HasPrefix(o, "/") {
			// clearly not a valid origin (looks like a path), skip and warn
			log.Printf("warning: skipping invalid allowed origin (appears to be a path): %q\n", o)
			continue
		}
		// If host:port was provided without scheme, assume http
		allowOrigins = append(allowOrigins, "http://"+o)
	}

	if len(allowOrigins) == 0 {
		// fallback to allowing all origins if nothing valid left
		allowOrigins = []string{"*"}
		log.Println("⚠️  No valid ALLOWED_ORIGINS found, defaulting to '*'")
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/api", func(c *gin.Context) {
		c.JSON(200, "")
	})

	v1 := r.Group("/dashboard/v1")
	{
		v1.POST("/auth/login", authHandler.Login)

		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			protected.GET("/payments", paymentHandler.ListPayments)
			protected.PUT("/payments/:id/review", paymentHandler.ReviewPayment)
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
