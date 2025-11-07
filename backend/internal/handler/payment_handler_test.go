package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"abasithdev.github.io/internal-cs-center-backend/internal/domain"
	"abasithdev.github.io/internal-cs-center-backend/internal/service"
	"abasithdev.github.io/internal-cs-center-backend/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func setupPaymentTest(t *testing.T) (*PaymentHandler, *gin.Engine, *storage.MemoryStore) {
	store := storage.NewMemoryStore()
	store.ClearPayments() // Start with clean slate

	paymentService := service.NewPaymentService(store)
	handler := NewPaymentHandler(paymentService)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/payments", handler.ListPayments)
	r.PUT("/payments/:id/review", handler.ReviewPayment)

	// Add test data
	now := time.Now()
	testPayments := []*domain.Payment{
		{
			ID:           "payment1",
			MerchantName: "Merchant A",
			Date:         now.Add(-time.Hour * 24),
			Amount:       100.0,
			Status:       "completed",
			Reviewed:     false,
		},
		{
			ID:           "payment2",
			MerchantName: "Merchant B",
			Date:         now,
			Amount:       50.0,
			Status:       "processing",
			Reviewed:     false,
		},
	}

	for _, p := range testPayments {
		store.UpdatePayment(p)
	}

	return handler, r, store
}

func TestPaymentHandler_ListPayments(t *testing.T) {
	_, r, _ := setupPaymentTest(t)

	tests := []struct {
		name         string
		queryParams  string
		wantCode     int
		wantItems    int
		wantStatus   string
		wantSortDesc bool
	}{
		{
			name:        "default params",
			queryParams: "",
			wantCode:    http.StatusOK,
			wantItems:   2,
		},
		{
			name:        "filter by status",
			queryParams: "status=completed",
			wantCode:    http.StatusOK,
			wantItems:   1,
			wantStatus:  "completed",
		},
		{
			name:        "search by merchant",
			queryParams: "search=payment1",
			wantCode:    http.StatusOK,
			wantItems:   1,
		},
		{
			name:        "pagination",
			queryParams: "page=1&size=1",
			wantCode:    http.StatusOK,
			wantItems:   1,
		},
		{
			name:        "sort by amount asc",
			queryParams: "sortBy=amount&orderBy=asc",
			wantCode:    http.StatusOK,
			wantItems:   2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/payments"
			if tt.queryParams != "" {
				url = fmt.Sprintf("/payments?%s", tt.queryParams)
			}

			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			require.Equal(t, tt.wantCode, w.Code)

			var resp map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			require.Contains(t, resp, "meta")
			require.Contains(t, resp, "summary")

			meta := resp["meta"].(map[string]interface{})
			data := meta["data"].([]interface{})

			if tt.wantItems > 0 {
				require.Len(t, data, tt.wantItems)
			}

			if tt.wantStatus != "" {
				payment := data[0].(map[string]interface{})
				require.Equal(t, tt.wantStatus, payment["status"])
			}
		})
	}
}

func TestPaymentHandler_ReviewPayment(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		role      string
		wantCode  int
		wantError bool
		errorMsg  string
	}{
		{
			name:      "unauthorized - cs role",
			id:        "payment1",
			role:      "cs",
			wantCode:  http.StatusUnauthorized,
			wantError: true,
			errorMsg:  "Forbidden",
		},
		{
			name:     "success - operational role",
			id:       "payment1",
			role:     "operational",
			wantCode: http.StatusOK,
		},
		{
			name:      "not found",
			id:        "nonexistent",
			role:      "operational",
			wantCode:  http.StatusNotFound,
			wantError: true,
			errorMsg:  "Payment not found",
		},
		{
			name:      "empty id",
			id:        "",
			role:      "operational",
			wantCode:  http.StatusBadRequest,
			wantError: true,
			errorMsg:  "id must not empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh router for each test
			gin.SetMode(gin.TestMode)
			r := gin.New()

			// Set up router middleware
			r.Use(func(c *gin.Context) {
				c.Set("role", tt.role)
			})

			handler, _, _ := setupPaymentTest(t)
			r.PUT("/payments/:id/review", handler.ReviewPayment)

			url := fmt.Sprintf("/payments/%s/review", tt.id)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPut, url, nil)
			r.ServeHTTP(w, req)

			require.Equal(t, tt.wantCode, w.Code)

			if tt.wantError {
				var resp map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				require.NoError(t, err)
				require.Contains(t, resp, "error")
				require.Equal(t, tt.errorMsg, resp["error"])
			}
		})
	}
}
