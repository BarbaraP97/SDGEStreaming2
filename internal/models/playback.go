// internal/models/playback.go
package models

import "time"

type PlaybackHistory struct {
	ID          int       `db:"id"`
	UserID      int       `db:"user_id"`
	ContentID   int       `db:"content_id"`
	ContentType string    `db:"content_type"`
	Progress    int       `db:"progress_seconds"`
	WatchedAt   time.Time `db:"watched_at"`
}

type Favorite struct {
	ID          int       `db:"id"`
	UserID      int       `db:"user_id"`
	ContentID   int       `db:"content_id"`
	ContentType string    `db:"content_type"`
	CreatedAt   time.Time `db:"created_at"`
}
