// internal/models/audiovisual.go
package models

// AudiovisualContent represents movies, series, or documentaries.
type AudiovisualContent struct {
	ID            int     `db:"id"`
	Title         string  `db:"title"`
	Type          string  `db:"type"`
	Genre         string  `db:"genre"`
	Duration      int     `db:"duration"` // minutes
	AgeRating     string  `db:"age_rating"`
	Synopsis      string  `db:"synopsis"`
	ReleaseYear   int     `db:"release_year"`
	Director      string  `db:"director"`
	AverageRating float64 `db:"average_rating"`
	IsAvailable   bool    `db:"is_available"`
}
