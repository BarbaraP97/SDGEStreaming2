package repositories

import (
	"SDGEStreaming/internal/db"
	"SDGEStreaming/internal/models"
	"database/sql"
	"fmt"
)

type FavoriteRepo interface {
	Create(f *models.Favorite) error
	Delete(userID, contentID int, contentType string) error
	FindByUserID(userID int) ([]models.Favorite, error)
}

type sqliteFavoriteRepo struct {
	conn *sql.DB
}

func NewFavoriteRepo() FavoriteRepo {
	return &sqliteFavoriteRepo{
		conn: db.GetDB(),
	}
}

func (r *sqliteFavoriteRepo) Create(f *models.Favorite) error {
	query := `
		INSERT INTO favorites (user_id, content_id, content_type)
		VALUES (?, ?, ?)
	`

	_, err := r.conn.Exec(query, f.UserID, f.ContentID, f.ContentType)
	if err != nil {
		return fmt.Errorf("error adding favorite: %w", err)
	}

	return nil
}

func (r *sqliteFavoriteRepo) Delete(userID, contentID int, contentType string) error {
	query := `
		DELETE FROM favorites
		WHERE user_id = ? AND content_id = ? AND content_type = ?
	`

	res, err := r.conn.Exec(query, userID, contentID, contentType)
	if err != nil {
		return fmt.Errorf("error deleting favorite: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("favorite not found")
	}

	return nil
}

func (r *sqliteFavoriteRepo) FindByUserID(userID int) ([]models.Favorite, error) {
	query := `
		SELECT id, user_id, content_id, content_type, added_at
		FROM favorites
		WHERE user_id = ?
		ORDER BY added_at DESC
	`

	rows, err := r.conn.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching favorites: %w", err)
	}
	defer rows.Close()

	var favorites []models.Favorite

	for rows.Next() {
		var f models.Favorite
		if err := rows.Scan(&f.ID, &f.UserID, &f.ContentID, &f.ContentType, &f.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning favorite row: %w", err)
		}
		favorites = append(favorites, f)
	}

	return favorites, nil
}
