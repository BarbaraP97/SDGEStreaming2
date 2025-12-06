package repositories

import (
	"SDGEStreaming/internal/db"
	"SDGEStreaming/internal/models"
	"database/sql"
	"fmt"
)

type PlaybackHistoryRepo interface {
	Create(history *models.PlaybackHistory) error
	UpdateProgress(userID, contentID int, contentType string, progress int) error
	FindByUserID(userID int) ([]models.PlaybackHistory, error)
	FindContinueWatching(userID int) ([]models.PlaybackHistory, error)
}

type sqlitePlaybackHistoryRepo struct {
	conn *sql.DB
}

func NewPlaybackHistoryRepo() PlaybackHistoryRepo {
	return &sqlitePlaybackHistoryRepo{
		conn: db.GetDB(),
	}
}

func (r *sqlitePlaybackHistoryRepo) Create(h *models.PlaybackHistory) error {
	query := `
		INSERT INTO playback_history (user_id, content_id, content_type, progress_seconds)
		VALUES (?, ?, ?, ?)
	`

	_, err := r.conn.Exec(query, h.UserID, h.ContentID, h.ContentType, h.Progress)
	if err != nil {
		return fmt.Errorf("error inserting playback history: %w", err)
	}

	return nil
}

func (r *sqlitePlaybackHistoryRepo) UpdateProgress(userID, contentID int, contentType string, progress int) error {
	query := `
		UPDATE playback_history
		SET progress_seconds = ?, watched_at = CURRENT_TIMESTAMP
		WHERE user_id = ? AND content_id = ? AND content_type = ?
	`

	res, err := r.conn.Exec(query, progress, userID, contentID, contentType)
	if err != nil {
		return fmt.Errorf("error updating playback progress: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no playback record found to update")
	}

	return nil
}

func (r *sqlitePlaybackHistoryRepo) FindByUserID(userID int) ([]models.PlaybackHistory, error) {
	query := `
		SELECT id, user_id, content_id, content_type, progress_seconds, watched_at
		FROM playback_history
		WHERE user_id = ?
		ORDER BY watched_at DESC
	`

	rows, err := r.conn.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching playback history: %w", err)
	}
	defer rows.Close()

	var history []models.PlaybackHistory

	for rows.Next() {
		var h models.PlaybackHistory
		if err := rows.Scan(&h.ID, &h.UserID, &h.ContentID, &h.ContentType, &h.Progress, &h.WatchedAt); err != nil {
			return nil, fmt.Errorf("error scanning playback history: %w", err)
		}
		history = append(history, h)
	}

	return history, nil
}

func (r *sqlitePlaybackHistoryRepo) FindContinueWatching(userID int) ([]models.PlaybackHistory, error) {
	query := `
		SELECT id, user_id, content_id, content_type, progress_seconds, watched_at
		FROM playback_history
		WHERE user_id = ?
		AND progress_seconds > 0
		ORDER BY watched_at DESC
		LIMIT 20
	`

	rows, err := r.conn.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching continue-watching list: %w", err)
	}
	defer rows.Close()

	var history []models.PlaybackHistory

	for rows.Next() {
		var h models.PlaybackHistory
		if err := rows.Scan(&h.ID, &h.UserID, &h.ContentID, &h.ContentType, &h.Progress, &h.WatchedAt); err != nil {
			return nil, fmt.Errorf("error scanning continue-watching rows: %w", err)
		}
		history = append(history, h)
	}

	return history, nil
}
