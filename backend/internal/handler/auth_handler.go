package handler

import (
	"net/http"

	"abasithdev.github.io/internal-cs-center-backend/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

func (auth *AuthHandler) Login(context *gin.Context) {
	var request loginRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := auth.auth.Authenticate(request.Email, request.Password)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credential"})
		return
	}

	token, err := auth.auth.GenerateToken(user)

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Failed generate token"})
		return
	}

	context.JSON(http.StatusOK, loginResponse{Token: token, Role: user.Role})
}
