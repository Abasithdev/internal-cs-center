package storage

import (
	"fmt"
	"testing"
	"time"

	"abasithdev.github.io/internal-cs-center-backend/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestMemoryStore_GetUserByEmail(t *testing.T) {
	store := NewMemoryStore()

	// Test existing user (from seed data)
	user, exists := store.GetUserByEmail("john-cs@durianpay.id")
	require.True(t, exists)
	require.Equal(t, "john-cs@durianpay.id", user.Email)
	require.Equal(t, "admin123", user.Password)
	require.Equal(t, "cs", user.Role)

	// Test non-existent user
	_, exists = store.GetUserByEmail("nonexistent@example.com")
	require.False(t, exists)
}

func TestMemoryStore_PaymentOperations(t *testing.T) {
	store := NewMemoryStore()
	store.ClearPayments() // Clear seeded payments
	now := time.Now()

	// Test GetPaymentList with empty store
	initial := store.GetPaymentList()
	require.Empty(t, initial)

	// Test UpdatePayment
	payment := &domain.Payment{
		ID:           "test1",
		MerchantName: "Test Merchant",
		Date:         now,
		Amount:       100.0,
		Status:       "processing",
		Reviewed:     false,
	}
	store.UpdatePayment(payment)

	// Test GetPaymentById
	retrieved, exists := store.GetPaymentById("test1")
	require.True(t, exists)
	require.Equal(t, payment.ID, retrieved.ID)
	require.Equal(t, payment.MerchantName, retrieved.MerchantName)
	require.Equal(t, payment.Amount, retrieved.Amount)
	require.Equal(t, payment.Status, retrieved.Status)
	require.Equal(t, payment.Reviewed, retrieved.Reviewed)

	// Test GetPaymentList after adding
	list := store.GetPaymentList()
	require.Len(t, list, 1)
	require.Equal(t, payment.ID, list[0].ID)

	// Test updating existing payment
	payment.Status = "completed"
	payment.Reviewed = true
	store.UpdatePayment(payment)

	updated, exists := store.GetPaymentById("test1")
	require.True(t, exists)
	require.Equal(t, "completed", updated.Status)
	require.True(t, updated.Reviewed)

	// Test non-existent payment
	_, exists = store.GetPaymentById("nonexistent")
	require.False(t, exists)
}

func TestMemoryStore_ConcurrentAccess(t *testing.T) {
	store := NewMemoryStore()
	store.ClearPayments() // Clear seeded payments
	done := make(chan bool)

	// Concurrent reads
	for i := 0; i < 10; i++ {
		go func() {
			store.GetUserByEmail("john-cs@durianpay.id")
			store.GetPaymentList()
			done <- true
		}()
	}

	// Concurrent writes
	for i := 0; i < 10; i++ {
		go func(i int) {
			payment := &domain.Payment{
				ID:     fmt.Sprintf("payment%d", i),
				Status: "processing",
			}
			store.UpdatePayment(payment)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 20; i++ {
		<-done
	}

	// Verify final state
	payments := store.GetPaymentList()
	require.Len(t, payments, 10)
}
