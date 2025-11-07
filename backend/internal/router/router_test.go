//go:build !skip_coverage

package router

import (
	"abasithdev.github.io/internal-cs-center-backend/internal/handler"
	"abasithdev.github.io/internal-cs-center-backend/internal/middleware"
	"abasithdev.github.io/internal-cs-center-backend/internal/service"
	"abasithdev.github.io/internal-cs-center-backend/internal/storage"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	store := storage.NewMemoryStore()
	authService := service.NewAuthService(store, []byte("donttellanyone"))
	paymentService := service.NewPaymentService(store)

	authHandler := handler.NewAuthHandler(authService)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	r := gin.Default()
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

	return r
}
