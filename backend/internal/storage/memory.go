package storage

import (
	"sync"

	"abasithdev.github.io/internal-cs-center-backend/internal/domain"
)

type MemoryStore struct {
	mu    sync.RWMutex
	users map[string]*domain.User
}

func NewMemoryStore() *MemoryStore {
	store := &MemoryStore{
		users: map[string]*domain.User{},
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
}

func (store *MemoryStore) GetUserByEmail(email string) (*domain.User, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	user, valid := store.users[email]
	return user, valid
}
