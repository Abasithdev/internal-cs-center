package middleware

import (
	"net/http"
	"strings"

	"abasithdev.github.io/internal-cs-center-backend/internal/service"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(auth *service.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		authPart := strings.SplitN(authHeader, " ", 2)
		if len(authPart) != 2 || strings.ToLower(authPart[0]) != "bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		token := authPart[1]
		claims, err := auth.ParseToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if role, ok := claims["role"].(string); ok {
			ctx.Set("role", role)
		}

		if email, ok := claims["email"].(string); ok {
			ctx.Set("email", email)
		}

		ctx.Next()
	}
}
