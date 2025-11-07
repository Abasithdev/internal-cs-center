package storage

/*
In memory storage
*/

import (
	"sync"
	"time"

	"abasithdev.github.io/internal-cs-center-backend/internal/domain"
	"github.com/google/uuid"
)

type MemoryStore struct {
	mu       sync.RWMutex
	users    map[string]*domain.User
	payments map[string]*domain.Payment
}

func NewMemoryStore() *MemoryStore {
	store := &MemoryStore{
		users:    map[string]*domain.User{},
		payments: map[string]*domain.Payment{},
	}

	store.seed()

	return store
}

func (store *MemoryStore) seed() {
	store.users["john-cs@durianpay.id"] = &domain.User{
		Email:    "john-cs@durianpay.id",
		Password: "admin123",
		Role:     "cs",
	}

	store.users["jane-operational@durianpay.id"] = &domain.User{
		Email:    "jane-operational@durianpay.id",
		Password: "admin123",
		Role:     "operational",
	}

	// seed for payments
	statuses := []string{"completed", "processing", "failed"}
	for i := 0; i < 20; i++ {
		id := uuid.New().String()
		payment := &domain.Payment{
			ID:           id,
			MerchantName: "Merchant" + id[:6],
			Date:         time.Now().Add(time.Duration(-i) * 24 * time.Hour),
			Amount:       float64(10000 + i*25),
			Status:       statuses[i%len(statuses)],
			Reviewed:     false,
		}

		store.payments[id] = payment
	}
}

// User
func (store *MemoryStore) GetUserByEmail(email string) (*domain.User, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	user, valid := store.users[email]
	return user, valid
}

// Payment
func (store *MemoryStore) GetPaymentList() []*domain.Payment {
	store.mu.RLock()
	defer store.mu.RUnlock()

	response := make([]*domain.Payment, 0, len(store.payments))

	for _, payment := range store.payments {
		response = append(response, payment)
	}

	return response
}

func (store *MemoryStore) GetPaymentById(id string) (*domain.Payment, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	payment, ok := store.payments[id]

	return payment, ok
}

func (store *MemoryStore) UpdatePayment(payment *domain.Payment) {
	store.mu.Lock()
	defer store.mu.Unlock()

	store.payments[payment.ID] = payment
}

func (store *MemoryStore) ClearPayments() {
	store.mu.Lock()
	defer store.mu.Unlock()

	store.payments = make(map[string]*domain.Payment)
}
