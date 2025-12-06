package repositories

import (
	"SDGEStreaming/internal/db"
	"SDGEStreaming/internal/models"
	"database/sql"
	"fmt"
	"time"
)

type SubscriptionRepo interface {
	Create(sub *models.Subscription) error
	UpdateUserPlan(userID int, planID int) error
	FindByUserID(userID int) (*models.Subscription, error)
	FindAll() ([]models.Subscription, error)
	Cancel(userID int) error
	GetPlanByID(planID int) (*models.Plan, error)
	GetAllPlans() ([]models.Plan, error)
}

type sqliteSubscriptionRepo struct {
	conn *sql.DB
}

func NewSubscriptionRepo() SubscriptionRepo {
	return &sqliteSubscriptionRepo{
		conn: db.GetDB(),
	}
}

//
// CREAR SUSCRIPCIÓN
//

func (r *sqliteSubscriptionRepo) Create(s *models.Subscription) error {
	query := `
		INSERT INTO subscriptions (user_id, plan_id, start_date, end_date, is_active)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := r.conn.Exec(query,
		s.UserID,
		s.PlanID,
		s.StartDate,
		s.EndDate,
		s.IsActive,
	)

	if err != nil {
		return fmt.Errorf("error creating subscription: %w", err)
	}
	return nil
}

//
// ACTUALIZAR PLAN DEL USUARIO
//

func (r *sqliteSubscriptionRepo) UpdateUserPlan(userID int, planID int) error {
	query := `
		UPDATE subscriptions
		SET plan_id = ?, start_date = ?, end_date = ?, is_active = 1
		WHERE user_id = ?
	`

	now := time.Now()
	end := now.AddDate(0, 1, 0) // un mes de suscripción

	_, err := r.conn.Exec(query,
		planID,
		now,
		end,
		userID,
	)

	if err != nil {
		return fmt.Errorf("error updating user plan: %w", err)
	}
	return nil
}

//
// OBTENER SUSCRIPCIÓN POR USUARIO
//

func (r *sqliteSubscriptionRepo) FindByUserID(userID int) (*models.Subscription, error) {
	query := `
		SELECT id, user_id, plan_id, start_date, end_date, is_active
		FROM subscriptions
		WHERE user_id = ?
	`

	row := r.conn.QueryRow(query, userID)

	var s models.Subscription
	err := row.Scan(
		&s.ID,
		&s.UserID,
		&s.PlanID,
		&s.StartDate,
		&s.EndDate,
		&s.IsActive,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("error scanning subscription: %w", err)
	}

	return &s, nil
}

// solucion del error GetPlanByID method defined
// add this method to the subscription_repo.go file
func (r *sqliteSubscriptionRepo) GetPlanByID(planID int) (*models.Plan, error) {
	query := `
		SELECT id, name, price
		FROM plans
		WHERE id = ?
	`

	row := r.conn.QueryRow(query, planID)

	var p models.Plan
	err := row.Scan(
		&p.ID,
		&p.Name,
		&p.Price,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("error scanning plan: %w", err)
	}

	return &p, nil
}

//
// LISTAR TODAS LAS SUSCRIPCIONES
//

func (r *sqliteSubscriptionRepo) FindAll() ([]models.Subscription, error) {
	query := `
		SELECT id, user_id, plan_id, start_date, end_date, is_active
		FROM subscriptions
		ORDER BY start_date DESC
	`

	rows, err := r.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error fetching subscriptions: %w", err)
	}
	defer rows.Close()

	var list []models.Subscription

	for rows.Next() {
		var s models.Subscription
		if err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.PlanID,
			&s.StartDate,
			&s.EndDate,
			&s.IsActive,
		); err != nil {
			return nil, fmt.Errorf("error scanning subscription row: %w", err)
		}
		list = append(list, s)
	}

	return list, nil
}

//
// CANCELAR SUSCRIPCIÓN
//

func (r *sqliteSubscriptionRepo) Cancel(userID int) error {
	query := `
		UPDATE subscriptions
		SET is_active = 0
		WHERE user_id = ?
	`

	_, err := r.conn.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("error canceling subscription: %w", err)
	}

	return nil
}

//
// OBTENER TODOS LOS PLANES
//

func (r *sqliteSubscriptionRepo) GetAllPlans() ([]models.Plan, error) {
	query := `
		SELECT id, name, price
		FROM plans
	`

	rows, err := r.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error fetching plans: %w", err)
	}
	defer rows.Close()

	var list []models.Plan

	for rows.Next() {
		var p models.Plan
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Price,
		); err != nil {
			return nil, fmt.Errorf("error scanning plan row: %w", err)
		}
		list = append(list, p)
	}

	return list, nil
}
