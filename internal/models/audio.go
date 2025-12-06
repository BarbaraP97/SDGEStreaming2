// internal/models/audio.go
package models

// AudioContent represents music, podcasts, or audiobooks.
type AudioContent struct {
	ID            int     `db:"id"`
	Title         string  `db:"title"`
	Type          string  `db:"type"`
	Genre         string  `db:"genre"`
	Duration      int     `db:"duration"` // minutes
	AgeRating     string  `db:"age_rating"`
	Artist        string  `db:"artist"`
	Album         string  `db:"album"`
	TrackNumber   int     `db:"track_number"`
	AverageRating float64 `db:"average_rating"`
	IsAvailable   bool    `db:"is_available"`
}
