// internal/models/playback.go
package models

import "time"

type PlaybackHistory struct {
	ID          int       `db:"id"`
	ProfileID   int       `db:"profile_id"`
	ContentID   int       `db:"content_id"`
	ContentType string    `db:"content_type"`
	Progress    int       `db:"progress_seconds"`
	IsCompleted bool      `db:"is_completed"`
	WatchedAt   time.Time `db:"watched_at"`
}

type Favorite struct {
	ID          int       `db:"id"`
	ProfileID   int       `db:"profile_id"`
	ContentID   int       `db:"content_id"`
	ContentType string    `db:"content_type"`
	CreatedAt   time.Time `db:"added_at"`
}