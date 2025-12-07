// internal/models/subscription.go
package models

import "time"

type Plan struct {
	ID    int     `db:"id"`
	Name  string  `db:"name"`
	Price float64 `db:"price"`
	MaxQuality string  `db:"max_quality"`
	MaxDevices int     `db:"max_devices"`
}

// Subscription representa la suscripción de un usuario a un plan.

type Subscription struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	PlanID    int       `db:"plan_id"`
	StartDate time.Time `db:"start_date"`
	EndDate   time.Time `db:"end_date"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type PaymentMethod struct {
	UserID         int       `db:"user_id"`
	CardNumber     string    `db:"card_number"`
	ExpirationDate string    `db:"expiration_date"`
	CVV            string    `db:"cvv"`
	CardHolder     string    `db:"card_holder_name"`
	Last4          string    `db:"card_number_last4"`
	ExpiryMonth    int       `db:"expiry_month"`
	ExpiryYear     int       `db:"expiry_year"`
	IsDefault      bool      `db:"is_default"`
	CreatedAt      time.Time `db:"created_at"`
}
