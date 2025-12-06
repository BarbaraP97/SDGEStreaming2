// internal/repositories/user_repo.go
package repositories

import (
	"SDGEStreaming/internal/db"
	"SDGEStreaming/internal/models"
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
	db := db.GetDB()

	query := `
		SELECT id, nombre, apellido, email, edad, password_hash
		FROM users
		ORDER BY id ASC
	`

	rows, err := db.Query(query)
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
			&u.PasswordHash,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

func (r *sqliteUserRepo) FindByID(id int) (*models.User, error) {
	db := db.GetDB()

	query := `
		SELECT id, nombre, apellido, email, edad, password_hash
		FROM users
		WHERE id = ?
	`

	var u models.User
	err := db.QueryRow(query, id).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Age,
		&u.PasswordHash,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *sqliteUserRepo) FindByEmail(email string) (*models.User, error) {
	db := db.GetDB()

	query := `
		SELECT id, nombre, apellido, email, edad, password_hash
		FROM users
		WHERE email = ?
	`

	var u models.User
	err := db.QueryRow(query, email).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Age,
		&u.PasswordHash,
	)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *sqliteUserRepo) Create(u *models.User) error {
	db := db.GetDB()

	query := `
		INSERT INTO users (nombre, apellido, email, edad, password_hash)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := db.Exec(query,
		u.Name,
		u.Email,
		u.Age,
		u.PasswordHash,
	)
	return err
}

func (r *sqliteUserRepo) Update(u *models.User) error {
	db := db.GetDB()

	query := `
		UPDATE users
		SET nombre = ?, apellido = ?, email = ?, edad = ?, password_hash = ?
		WHERE id = ?
	`

	_, err := db.Exec(query,
		u.Name,
		u.Email,
		u.Age,
		u.PasswordHash,
		u.ID,
	)
	return err
}

func (r *sqliteUserRepo) Delete(id int) error {
	db := db.GetDB()

	query := `
		DELETE FROM users WHERE id = ?
	`

	_, err := db.Exec(query, id)
	return err
}

func (r *sqliteUserRepo) UpdatePlan(userID int, planID int) error {
	db := db.GetDB()

	query := `
		UPDATE users
		SET plan_id = ?
		WHERE id = ?
	`

	_, err := db.Exec(query, planID, userID)
	return err
}

func (r *sqliteUserRepo) GetDefaultPaymentMethod(userID int) (*models.PaymentMethod, error) {
	db := db.GetDB()

	query := `
		SELECT user_id, card_number, expiration_date, cvv
		FROM payment_methods
		WHERE user_id = ?
		LIMIT 1
	`

	var pm models.PaymentMethod
	err := db.QueryRow(query, userID).Scan(
		&pm.UserID,
		&pm.CardNumber,
		&pm.ExpirationDate,
		&pm.CVV,
	)
	if err != nil {
		return nil, err
	}

	return &pm, nil
}
func (r *sqliteUserRepo) AddPaymentMethod(pm *models.PaymentMethod) error {
	db := db.GetDB()

	query := `
		INSERT INTO payment_methods (user_id, card_number, expiration_date, cvv)
		VALUES (?, ?, ?, ?)
	`

	_, err := db.Exec(query,
		pm.UserID,
		pm.CardNumber,
		pm.ExpirationDate,
		pm.CVV,
	)
	return err
}
