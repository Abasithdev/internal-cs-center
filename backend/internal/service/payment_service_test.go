package service

import (
	"testing"
	"time"

	"abasithdev.github.io/internal-cs-center-backend/internal/domain"
	"abasithdev.github.io/internal-cs-center-backend/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestPaymentService_GetList(t *testing.T) {
	// Create a memory store with test data
	store := storage.NewMemoryStore()
	store.ClearPayments() // Clear seeded payments
	now := time.Now()

	testPayments := []*domain.Payment{
		{
			ID:           "payment1",
			MerchantName: "Merchant A",
			Date:         now.Add(-time.Hour * 24), // 1 day ago
			Amount:       100.0,
			Status:       "completed",
			Reviewed:     true,
		},
		{
			ID:           "payment2",
			MerchantName: "Merchant B",
			Date:         now,
			Amount:       50.0,
			Status:       "processing",
			Reviewed:     false,
		},
		{
			ID:           "payment3",
			MerchantName: "Merchant C",
			Date:         now.Add(-time.Hour * 48), // 2 days ago
			Amount:       75.0,
			Status:       "failed",
			Reviewed:     false,
		},
	}

	// Add test payments to store
	for _, p := range testPayments {
		store.UpdatePayment(p)
	}

	service := NewPaymentService(store)

	tests := []struct {
		name     string
		request  ListRequest
		want     int    // expected items in result
		wantSize int    // expected page size
		wantPage int    // expected page number
		sortBy   string // field to sort by
		orderBy  string // asc or desc
	}{
		{
			name: "default pagination",
			request: ListRequest{
				Page: 1,
				Size: 10,
			},
			want:     3,
			wantSize: 10,
			wantPage: 1,
		},
		{
			name: "filter by status",
			request: ListRequest{
				Status: "completed",
				Page:   1,
				Size:   10,
			},
			want:     1,
			wantSize: 10,
			wantPage: 1,
		},
		{
			name: "search by ID",
			request: ListRequest{
				Search: "payment1",
				Page:   1,
				Size:   10,
			},
			want:     1,
			wantSize: 10,
			wantPage: 1,
		},
		{
			name: "sort by amount asc",
			request: ListRequest{
				Page:    1,
				Size:    10,
				SortBy:  "amount",
				OrderBy: "asc",
			},
			want:     3,
			wantSize: 10,
			wantPage: 1,
		},
		{
			name: "sort by date desc",
			request: ListRequest{
				Page:    1,
				Size:    10,
				SortBy:  "date",
				OrderBy: "desc",
			},
			want:     3,
			wantSize: 10,
			wantPage: 1,
		},
		{
			name: "pagination - page 1 size 2",
			request: ListRequest{
				Page: 1,
				Size: 2,
			},
			want:     2,
			wantSize: 2,
			wantPage: 1,
		},
		{
			name: "pagination - page 2 size 2",
			request: ListRequest{
				Page: 2,
				Size: 2,
			},
			want:     1,
			wantSize: 2,
			wantPage: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.GetList(tt.request)

			require.Equal(t, tt.want, len(result.Data), "expected %d items, got %d", tt.want, len(result.Data))
			require.Equal(t, tt.wantSize, result.Size, "expected page size %d, got %d", tt.wantSize, result.Size)
			require.Equal(t, tt.wantPage, result.Page, "expected page number %d, got %d", tt.wantPage, result.Page)

			if tt.sortBy == "amount" && tt.orderBy == "asc" {
				// Verify ascending sort by amount
				for i := 1; i < len(result.Data); i++ {
					require.GreaterOrEqual(t, result.Data[i].Amount, result.Data[i-1].Amount,
						"amounts should be in ascending order")
				}
			}

			if tt.sortBy == "date" && tt.orderBy == "desc" {
				// Verify descending sort by date
				for i := 1; i < len(result.Data); i++ {
					require.True(t, result.Data[i].Date.Before(result.Data[i-1].Date) || result.Data[i].Date.Equal(result.Data[i-1].Date),
						"dates should be in descending order")
				}
			}
		})
	}
}

func TestPaymentService_GetStatusSummary(t *testing.T) {
	store := storage.NewMemoryStore()
	store.ClearPayments() // Clear seeded payments
	now := time.Now()

	testPayments := []*domain.Payment{
		{ID: "1", Status: "completed", Date: now},
		{ID: "2", Status: "completed", Date: now},
		{ID: "3", Status: "processing", Date: now},
		{ID: "4", Status: "failed", Date: now},
		{ID: "5", Status: "completed", Date: now},
	}

	for _, p := range testPayments {
		store.UpdatePayment(p)
	}

	service := NewPaymentService(store)
	completed, processing, failed := service.GetStatusSummary()

	require.Equal(t, 3, completed, "expected 3 completed payments")
	require.Equal(t, 1, processing, "expected 1 processing payment")
	require.Equal(t, 1, failed, "expected 1 failed payment")
}

func TestPaymentService_Review(t *testing.T) {
	store := storage.NewMemoryStore()
	store.ClearPayments() // Clear seeded payments
	service := NewPaymentService(store)

	// Test successful review
	payment := &domain.Payment{
		ID:       "test1",
		Status:   "processing",
		Reviewed: false,
	}
	store.UpdatePayment(payment)

	err := service.Review("test1")
	require.NoError(t, err)

	updated, exists := store.GetPaymentById("test1")
	require.True(t, exists)
	require.True(t, updated.Reviewed)

	// Test review non-existent payment
	err = service.Review("nonexistent")
	require.Error(t, err)
	require.Contains(t, err.Error(), "paymentId: nonexistent")
}

func TestPaymentService_GetTotalByFilter(t *testing.T) {
	store := storage.NewMemoryStore()
	store.ClearPayments() // Clear seeded payments
	now := time.Now()

	testPayments := []*domain.Payment{
		{ID: "1", Status: "completed", Date: now},
		{ID: "2", Status: "completed", Date: now},
		{ID: "3", Status: "processing", Date: now},
		{ID: "4", Status: "failed", Date: now},
	}

	for _, p := range testPayments {
		store.UpdatePayment(p)
	}

	service := NewPaymentService(store)

	tests := []struct {
		name    string
		request ListRequest
		want    int
	}{
		{
			name:    "no filter",
			request: ListRequest{},
			want:    4,
		},
		{
			name: "filter by completed status",
			request: ListRequest{
				Status: "completed",
			},
			want: 2,
		},
		{
			name: "filter by ID search",
			request: ListRequest{
				Search: "1",
			},
			want: 1,
		},
		{
			name: "no matches",
			request: ListRequest{
				Search: "nonexistent",
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.GetTotalByFilter(tt.request)
			require.Equal(t, tt.want, got)
		})
	}
}
