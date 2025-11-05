package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"abasithdev.github.io/internal-cs-center-backend/internal/router"
	"github.com/stretchr/testify/require"
)

type loginResp struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

func TestLoginSuccess(t *testing.T) {
	r := router.NewRouter()

	payload := map[string]string{
		"email":    "john-cs@durianpay.id",
		"password": "admin123",
	}

	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/dashboard/v1/auth/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp loginResp
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.NotEmpty(t, resp.Token)
	require.Equal(t, "cs", resp.Role)
}

func TestLoginFailed(t *testing.T) {
	r := router.NewRouter()

	payload := map[string]string{
		"email":    "john-cs@durianpay.id",
		"password": "wrongpassword",
	}

	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/dashboard/v1/auth/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestPaymentsUnauthorized(t *testing.T) {
	r := router.NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/dashboard/v1/payments", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestPaymentsAuthorized(t *testing.T) {
	r := router.NewRouter()

	// login first to get a token
	payload := map[string]string{
		"email":    "john-cs@durianpay.id",
		"password": "admin123",
	}

	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/dashboard/v1/auth/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var lr loginResp
	err := json.Unmarshal(w.Body.Bytes(), &lr)
	require.NoError(t, err)
	require.NotEmpty(t, lr.Token)

	// call protected endpoint with Bearer token
	req2 := httptest.NewRequest(http.MethodGet, "/dashboard/v1/payments", nil)
	req2.Header.Set("Authorization", "Bearer "+lr.Token)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	require.Equal(t, http.StatusOK, w2.Code)

	var respStr string
	err = json.Unmarshal(w2.Body.Bytes(), &respStr)
	require.NoError(t, err)
	require.Equal(t, "authorized", respStr)
}
