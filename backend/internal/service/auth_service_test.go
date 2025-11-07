package service

import (
	"testing"
	"time"

	"abasithdev.github.io/internal-cs-center-backend/internal/domain"
	"abasithdev.github.io/internal-cs-center-backend/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestAuthService_Authenticate(t *testing.T) {
	store := storage.NewMemoryStore() // Uses seed data with test users
	secret := []byte("test-secret-key")
	service := NewAuthService(store, secret)

	tests := []struct {
		name      string
		email     string
		password  string
		wantError bool
		wantRole  string
	}{
		{
			name:      "valid cs user",
			email:     "john-cs@durianpay.id",
			password:  "admin123",
			wantError: false,
			wantRole:  "cs",
		},
		{
			name:      "valid operational user",
			email:     "jane-operational@durianpay.id",
			password:  "admin123",
			wantError: false,
			wantRole:  "operational",
		},
		{
			name:      "wrong password",
			email:     "john-cs@durianpay.id",
			password:  "wrongpass",
			wantError: true,
		},
		{
			name:      "non-existent user",
			email:     "nobody@example.com",
			password:  "anypass",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := service.Authenticate(tt.email, tt.password)

			if tt.wantError {
				require.Error(t, err)
				require.Nil(t, user)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, user)
			require.Equal(t, tt.email, user.Email)
			require.Equal(t, tt.wantRole, user.Role)
		})
	}
}

func TestAuthService_GenerateAndParseToken(t *testing.T) {
	store := storage.NewMemoryStore()
	secret := []byte("test-secret-key")
	service := NewAuthService(store, secret)

	// Test valid token generation and parsing
	user := &domain.User{
		Email: "test@example.com",
		Role:  "cs",
	}

	token, err := service.GenerateToken(user)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// Parse and verify the token
	claims, err := service.ParseToken(token)
	require.NoError(t, err)
	require.NotNil(t, claims)

	// Verify claims
	require.Equal(t, user.Email, claims["email"])
	require.Equal(t, user.Role, claims["role"])
	require.NotZero(t, claims["exp"])

	// Test expired token
	expiredService := NewAuthService(store, secret)
	expiredService.tokenValidation = -time.Hour // Token that expired 1 hour ago

	expiredToken, _ := expiredService.GenerateToken(user)
	_, err = service.ParseToken(expiredToken)
	require.Error(t, err)
	require.Contains(t, err.Error(), "token is expired")

	// Test invalid token format
	_, err = service.ParseToken("invalid.token.format")
	require.Error(t, err)

	// Test token with invalid signature
	_, err = service.ParseToken(token + "modified")
	require.Error(t, err)

	// Test token signed with different key
	differentService := NewAuthService(store, []byte("different-secret"))
	differentToken, _ := differentService.GenerateToken(user)
	_, err = service.ParseToken(differentToken)
	require.Error(t, err)
}

func TestAuthService_TokenValidation(t *testing.T) {
	store := storage.NewMemoryStore()
	secret := []byte("test-secret-key")
	service := NewAuthService(store, secret)

	user := &domain.User{
		Email: "test@example.com",
		Role:  "cs",
	}

	// Test default validation period
	token, err := service.GenerateToken(user)
	require.NoError(t, err)

	claims, err := service.ParseToken(token)
	require.NoError(t, err)

	exp, ok := claims["exp"].(float64)
	require.True(t, ok)

	// Verify expiration is about 24 hours in the future (within 1 second tolerance)
	expectedExp := float64(time.Now().Add(24 * time.Hour).Unix())
	require.InDelta(t, expectedExp, exp, 1.0)

	// Test custom validation period
	customService := NewAuthService(store, secret)
	customService.tokenValidation = time.Hour // 1 hour

	token, err = customService.GenerateToken(user)
	require.NoError(t, err)

	claims, err = customService.ParseToken(token)
	require.NoError(t, err)

	exp, ok = claims["exp"].(float64)
	require.True(t, ok)

	// Verify expiration is about 1 hour in the future
	expectedExp = float64(time.Now().Add(time.Hour).Unix())
	require.InDelta(t, expectedExp, exp, 1.0)
}
