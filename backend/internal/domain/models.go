package domain

import "time"

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Payment struct {
	ID           string    `json:"id"`
	MerchantName string    `json:"merchant_name"`
	Date         time.Time `json:"date"`
	Amount       float64   `json:"amount"`
	Status       string    `json:"status"`
	Reviewed     bool      `json:"reviewed"`
}
