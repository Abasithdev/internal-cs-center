package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"abasithdev.github.io/internal-cs-center-backend/internal/service"
	"abasithdev.github.io/internal-cs-center-backend/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func setupAuthTest(t *testing.T) (*AuthHandler, *gin.Engine) {
	store := storage.NewMemoryStore()
	authService := service.NewAuthService(store, []byte("donttellanyone"))
	handler := NewAuthHandler(authService)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/login", handler.Login)

	return handler, r
}

func TestAuthHandler_Login(t *testing.T) {
	_, r := setupAuthTest(t)

	tests := []struct {
		name      string
		request   loginRequest
		wantCode  int
		wantError bool
		wantToken bool
		wantRole  string
		errorMsg  string
	}{
		{
			name: "valid cs user",
			request: loginRequest{
				Email:    "john-cs@durianpay.id",
				Password: "admin123",
			},
			wantCode:  http.StatusOK,
			wantToken: true,
			wantRole:  "cs",
		},
		{
			name: "valid operational user",
			request: loginRequest{
				Email:    "jane-operational@durianpay.id",
				Password: "admin123",
			},
			wantCode:  http.StatusOK,
			wantToken: true,
			wantRole:  "operational",
		},
		{
			name: "invalid password",
			request: loginRequest{
				Email:    "john-cs@durianpay.id",
				Password: "wrong",
			},
			wantCode:  http.StatusUnauthorized,
			wantError: true,
			errorMsg:  "Invalid credential",
		},
		{
			name: "invalid email format",
			request: loginRequest{
				Email:    "invalid-email",
				Password: "admin123",
			},
			wantCode:  http.StatusBadRequest,
			wantError: true,
		},
		{
			name: "missing password",
			request: loginRequest{
				Email: "john-cs@durianpay.id",
			},
			wantCode:  http.StatusBadRequest,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := json.Marshal(tt.request)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(b))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			require.Equal(t, tt.wantCode, w.Code)

			if tt.wantError {
				var resp map[string]string
				err = json.Unmarshal(w.Body.Bytes(), &resp)
				require.NoError(t, err)
				require.Contains(t, resp, "error")
				if tt.errorMsg != "" {
					require.Equal(t, tt.errorMsg, resp["error"])
				}
				return
			}

			var resp loginResponse
			err = json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)
			if tt.wantToken {
				require.NotEmpty(t, resp.Token)
			}
			if tt.wantRole != "" {
				require.Equal(t, tt.wantRole, resp.Role)
			}
		})
	}
}
