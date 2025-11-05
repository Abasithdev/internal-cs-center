package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"abasithdev.github.io/internal-cs-center-backend/internal/domain"
	"abasithdev.github.io/internal-cs-center-backend/internal/storage"
)

type AuthService struct {
	jwtSecret       []byte
	tokenValidation time.Duration
	store           *storage.MemoryStore
}

func NewAuthService(store *storage.MemoryStore, secret []byte) *AuthService {
	return &AuthService{jwtSecret: secret, tokenValidation: time.Hour * 24, store: store}
}

func (auth *AuthService) Authenticate(email, password string) (*domain.User, error) {
	user, valid := auth.store.GetUserByEmail(email)

	if !valid {
		return nil, errors.New("Invalid user")
	}

	if user.Password != password {
		return nil, errors.New("Invalid user")
	}

	return user, nil
}

func (auth *AuthService) GenerateToken(user *domain.User) (string, error) {
	claim := jwt.MapClaims{
		"email": user.Email,
		"role":  user.Role,
		// use standard exp claim (unix timestamp)
		"exp": time.Now().Add(auth.tokenValidation).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	// sign using the configured secret (not the token object)
	return token.SignedString(auth.jwtSecret)
}

func (auth *AuthService) ParseToken(tokenStr string) (jwt.MapClaims, error) {
	parsed, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		// enforce HMAC signing method
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return auth.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claim, ok := parsed.Claims.(jwt.MapClaims); ok && parsed.Valid {
		return claim, nil
	}

	return nil, errors.New("Invalid token")
}
