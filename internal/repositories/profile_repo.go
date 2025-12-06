// internal/repositories/profile_repo.go
package repositories

import (
	"SDGEStreaming/internal/db"
	"SDGEStreaming/internal/models"
	"database/sql"
)

type ProfileRepo interface {
	Create(p *models.Profile) error
	FindByID(id int) (*models.Profile, error)
	FindByUserID(userID int) ([]models.Profile, error)
	Update(p *models.Profile) error
	Delete(id int) error
	CountByUserID(userID int) (int, error)
}

type sqliteProfileRepo struct{}

func NewProfileRepo() ProfileRepo {
	return &sqliteProfileRepo{}
}

func (r *sqliteProfileRepo) Create(p *models.Profile) error {
	conn := db.GetDB()
	query := `
		INSERT INTO profiles (user_id, name, type, age_rating, avatar, is_main, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	result, err := conn.Exec(query, p.UserID, p.Name, p.Type, p.AgeRating, p.Avatar, p.IsMain, p.CreatedAt)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	p.ID = int(id)
	return nil
}

func (r *sqliteProfileRepo) FindByID(id int) (*models.Profile, error) {
	conn := db.GetDB()
	query := `
		SELECT id, user_id, name, type, age_rating, avatar, is_main, created_at
		FROM profiles WHERE id = ?
	`
	var p models.Profile
	err := conn.QueryRow(query, id).Scan(&p.ID, &p.UserID, &p.Name, &p.Type, &p.AgeRating, &p.Avatar, &p.IsMain, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &p, err
}

func (r *sqliteProfileRepo) FindByUserID(userID int) ([]models.Profile, error) {
	conn := db.GetDB()
	query := `
		SELECT id, user_id, name, type, age_rating, avatar, is_main, created_at
		FROM profiles WHERE user_id = ? ORDER BY is_main DESC, created_at ASC
	`
	rows, err := conn.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []models.Profile
	for rows.Next() {
		var p models.Profile
		if err := rows.Scan(&p.ID, &p.UserID, &p.Name, &p.Type, &p.AgeRating, &p.Avatar, &p.IsMain, &p.CreatedAt); err != nil {
			return nil, err
		}
		profiles = append(profiles, p)
	}
	return profiles, nil
}

func (r *sqliteProfileRepo) Update(p *models.Profile) error {
	conn := db.GetDB()
	query := `
		UPDATE profiles SET name = ?, type = ?, age_rating = ?, avatar = ?
		WHERE id = ?
	`
	_, err := conn.Exec(query, p.Name, p.Type, p.AgeRating, p.Avatar, p.ID)
	return err
}

func (r *sqliteProfileRepo) Delete(id int) error {
	conn := db.GetDB()
	_, err := conn.Exec("DELETE FROM profiles WHERE id = ?", id)
	return err
}

func (r *sqliteProfileRepo) CountByUserID(userID int) (int, error) {
	conn := db.GetDB()
	var count int
	err := conn.QueryRow("SELECT COUNT(*) FROM profiles WHERE user_id = ?", userID).Scan(&count)
	return count, err
}