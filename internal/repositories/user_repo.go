// internal/repositories/user_repo.go
package repositories

import (
	"SDGEStreaming/internal/db"
	"SDGEStreaming/internal/models"
	"database/sql"
	"fmt"
)

type UserRepo interface {
	FindAll() ([]models.User, error)
	FindByID(id int) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Create(u *models.User) error
	Update(u *models.User) error
	Delete(id int) error
	UpdatePlan(userID int, planID int) error
	AddPaymentMethod(pm *models.PaymentMethod) error
	GetDefaultPaymentMethod(userID int) (*models.PaymentMethod, error)
}

type sqliteUserRepo struct{}

func NewUserRepo() UserRepo {
	return &sqliteUserRepo{}
}

func (r *sqliteUserRepo) FindAll() ([]models.User, error) {
	conn := db.GetDB()

	query := `
		SELECT id, name, email, age, plan_id, age_rating, is_admin, password_hash, created_at, last_login
		FROM users
		ORDER BY id ASC
	`

	rows, err := conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var u models.User
		err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
			&u.Age,
			&u.PlanID,
			&u.AgeRating,
			&u.IsAdmin,
			&u.PasswordHash,
			&u.CreatedAt,
			&u.LastLogin,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (r *sqliteUserRepo) FindByID(id int) (*models.User, error) {
	conn := db.GetDB()

	query := `
		SELECT id, name, email, age, plan_id, age_rating, is_admin, password_hash, created_at, last_login
		FROM users
		WHERE id = ?
	`

	var u models.User
	err := conn.QueryRow(query, id).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Age,
		&u.PlanID,
		&u.AgeRating,
		&u.IsAdmin,
		&u.PasswordHash,
		&u.CreatedAt,
		&u.LastLogin,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *sqliteUserRepo) FindByEmail(email string) (*models.User, error) {
	conn := db.GetDB()

	query := `
		SELECT id, name, email, age, plan_id, age_rating, is_admin, password_hash, created_at, last_login
		FROM users
		WHERE email = ?
	`

	var u models.User
	err := conn.QueryRow(query, email).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Age,
		&u.PlanID,
		&u.AgeRating,
		&u.IsAdmin,
		&u.PasswordHash,
		&u.CreatedAt,
		&u.LastLogin,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *sqliteUserRepo) Create(u *models.User) error {
	conn := db.GetDB()

	query := `
		INSERT INTO users (name, email, age, plan_id, age_rating, is_admin, password_hash, created_at, last_login)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := conn.Exec(query,
		u.Name,
		u.Email,
		u.Age,
		u.PlanID,
		u.AgeRating,
		u.IsAdmin,
		u.PasswordHash,
		u.CreatedAt,
		u.LastLogin,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = int(id)
	return nil
}

func (r *sqliteUserRepo) Update(u *models.User) error {
	conn := db.GetDB()

	query := `
		UPDATE users
		SET name = ?, email = ?, age = ?, plan_id = ?, age_rating = ?, is_admin = ?, password_hash = ?
		WHERE id = ?
	`

	_, err := conn.Exec(query,
		u.Name,
		u.Email,
		u.Age,
		u.PlanID,
		u.AgeRating,
		u.IsAdmin,
		u.PasswordHash,
		u.ID,
	)
	return err
}

func (r *sqliteUserRepo) Delete(id int) error {
	conn := db.GetDB()

	query := `
		DELETE FROM users WHERE id = ?
	`

	_, err := conn.Exec(query, id)
	return err
}

func (r *sqliteUserRepo) UpdatePlan(userID int, planID int) error {
	conn := db.GetDB()

	query := `
		UPDATE users
		SET plan_id = ?
		WHERE id = ?
	`

	_, err := conn.Exec(query, planID, userID)
	return err
}

func (r *sqliteUserRepo) GetDefaultPaymentMethod(userID int) (*models.PaymentMethod, error) {
	conn := db.GetDB()

	query := `
		SELECT user_id, card_holder_name, card_number_last4, expiry_month, expiry_year, is_default, created_at
		FROM payment_methods
		WHERE user_id = ? AND is_default = 1
		LIMIT 1
	`

	var pm models.PaymentMethod
	err := conn.QueryRow(query, userID).Scan(
		&pm.UserID,
		&pm.CardHolder,
		&pm.Last4,
		&pm.ExpiryMonth,
		&pm.ExpiryYear,
		&pm.IsDefault,
		&pm.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no se encontró método de pago")
	}
	if err != nil {
		return nil, err
	}

	return &pm, nil
}

func (r *sqliteUserRepo) AddPaymentMethod(pm *models.PaymentMethod) error {
	conn := db.GetDB()

	query := `
		INSERT INTO payment_methods (
			user_id,
			card_holder_name,
			card_number,
			expiration_date,
			cvv,
			card_number_last4,
			expiry_month,
			expiry_year,
			is_default
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := conn.Exec(query,
		pm.UserID,
		pm.CardHolder,
		pm.CardNumber,
		pm.ExpirationDate,
		pm.CVV,
		pm.Last4,
		pm.ExpiryMonth,
		pm.ExpiryYear,
		pm.IsDefault,
	)
	return err
}
