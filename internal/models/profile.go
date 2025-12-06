// internal/models/profile.go
package models

import "time"

type Profile struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Name      string    `db:"name"`
	Type      string    `db:"type"` // kids, teen, adult
	AgeRating string    `db:"age_rating"`
	Avatar    string    `db:"avatar"`
	IsMain    bool      `db:"is_main"`
	CreatedAt time.Time `db:"created_at"`
}

type Playlist struct {
	ID          int       `db:"id"`
	ProfileID   int       `db:"profile_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
}

type PlaylistItem struct {
	ID          int    `db:"id"`
	PlaylistID  int    `db:"playlist_id"`
	ContentID   int    `db:"content_id"`
	ContentType string `db:"content_type"`
	Position    int    `db:"position"`
}