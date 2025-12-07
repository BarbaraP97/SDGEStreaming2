// internal/models/user.go
package models

import "time"

// User represents a registered user in the system.
type User struct {
	ID           int       `db:"id"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	Age          int       `db:"age"`
	PlanID       int       `db:"plan_id"`
	AgeRating    string    `db:"age_rating"`
	IsAdmin      bool      `db:"is_admin"`
	CreatedAt    time.Time `db:"created_at"`
	LastLogin    time.Time `db:"last_login"`
	PasswordHash string    `db:"password_hash"`
}
