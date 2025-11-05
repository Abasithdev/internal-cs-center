package router

import (
	"abasithdev.github.io/internal-cs-center-backend/internal/handler"
	"abasithdev.github.io/internal-cs-center-backend/internal/middleware"
	"abasithdev.github.io/internal-cs-center-backend/internal/service"
	"abasithdev.github.io/internal-cs-center-backend/internal/storage"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter() *gin.Engine {
	store := storage.NewMemoryStore()
	authService := service.NewAuthService(store, []byte("donttellanyone"))

	authHandler := handler.NewAuthHandler(authService)

	r := gin.Default()

	r.GET("/api", func(c *gin.Context) {
		c.JSON(200, "")
	})

	v1 := r.Group("/dashboard/v1")
	{
		v1.POST("/auth/login", authHandler.Login)

		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			protected.GET("/payments", func(ctx *gin.Context) {
				ctx.JSON(200, "authorized")
			})
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
