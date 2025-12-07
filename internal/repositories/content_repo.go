// internal/repositories/content_repo.go
package repositories

import (
	"SDGEStreaming/internal/db"
	"SDGEStreaming/internal/models"
	"database/sql"
	"fmt"
)

type ContentRepo interface {
	// Audiovisual
	CreateAudiovisual(content *models.AudiovisualContent) error
	FindAudiovisualByID(id int) (*models.AudiovisualContent, error)
	FindAllAudiovisual() ([]models.AudiovisualContent, error)
	SearchAudiovisualByTitle(title string) ([]models.AudiovisualContent, error)

	// Audio
	CreateAudio(content *models.AudioContent) error
	FindAudioByID(id int) (*models.AudioContent, error)
	FindAllAudio() ([]models.AudioContent, error)
	SearchAudioByTitle(title string) ([]models.AudioContent, error)

	// Ratings
	UpdateAverageRating(contentID int, contentType string, avg float64) error

	// Filtrado por edad
	FindAllAudiovisualAllowed(userAgeRating string) ([]models.AudiovisualContent, error)
	FindAllAudioAllowed(userAgeRating string) ([]models.AudioContent, error)
}

type sqliteContentRepo struct{}

func NewContentRepo() ContentRepo {
	return &sqliteContentRepo{}
}

// --- AUDIOVISUAL ---

func (r *sqliteContentRepo) CreateAudiovisual(content *models.AudiovisualContent) error {
	conn := db.GetDB()

	query := `
		INSERT INTO audiovisual_content (title, type, genre, duration, age_rating, synopsis, release_year, director, average_rating, is_available)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := conn.Exec(query,
		content.Title,
		content.Type,
		content.Genre,
		content.Duration,
		content.AgeRating,
		content.Synopsis,
		content.ReleaseYear,
		content.Director,
		content.AverageRating,
		content.IsAvailable,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	content.ID = int(id)
	return nil
}

func (r *sqliteContentRepo) FindAudiovisualByID(id int) (*models.AudiovisualContent, error) {
	conn := db.GetDB()

	query := `
		SELECT id, title, type, genre, duration, age_rating, synopsis, release_year, director, average_rating, is_available
		FROM audiovisual_content
		WHERE id = ?
	`

	var c models.AudiovisualContent
	err := conn.QueryRow(query, id).Scan(
		&c.ID,
		&c.Title,
		&c.Type,
		&c.Genre,
		&c.Duration,
		&c.AgeRating,
		&c.Synopsis,
		&c.ReleaseYear,
		&c.Director,
		&c.AverageRating,
		&c.IsAvailable,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("contenido audiovisual no encontrado")
	}
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *sqliteContentRepo) FindAllAudiovisual() ([]models.AudiovisualContent, error) {
	conn := db.GetDB()

	query := `
		SELECT id, title, type, genre, duration, age_rating, synopsis, release_year, director, average_rating, is_available
		FROM audiovisual_content
		WHERE is_available = 1
		ORDER BY average_rating DESC
	`

	rows, err := conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contents []models.AudiovisualContent
	for rows.Next() {
		var c models.AudiovisualContent
		err := rows.Scan(
			&c.ID,
			&c.Title,
			&c.Type,
			&c.Genre,
			&c.Duration,
			&c.AgeRating,
			&c.Synopsis,
			&c.ReleaseYear,
			&c.Director,
			&c.AverageRating,
			&c.IsAvailable,
		)
		if err != nil {
			return nil, err
		}
		contents = append(contents, c)
	}

	return contents, nil
}

func (r *sqliteContentRepo) SearchAudiovisualByTitle(title string) ([]models.AudiovisualContent, error) {
	conn := db.GetDB()

	query := `
		SELECT id, title, type, genre, duration, age_rating, synopsis, release_year, director, average_rating, is_available
		FROM audiovisual_content
		WHERE title LIKE ? AND is_available = 1
		ORDER BY average_rating DESC
	`

	rows, err := conn.Query(query, "%"+title+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contents []models.AudiovisualContent
	for rows.Next() {
		var c models.AudiovisualContent
		err := rows.Scan(
			&c.ID,
			&c.Title,
			&c.Type,
			&c.Genre,
			&c.Duration,
			&c.AgeRating,
			&c.Synopsis,
			&c.ReleaseYear,
			&c.Director,
			&c.AverageRating,
			&c.IsAvailable,
		)
		if err != nil {
			return nil, err
		}
		contents = append(contents, c)
	}

	return contents, nil
}

func (r *sqliteContentRepo) FindAllAudiovisualAllowed(userAgeRating string) ([]models.AudiovisualContent, error) {
	conn := db.GetDB()

	var query string

	// Regla simple:
	// - Niño: solo G
	// - Adolescente: G, PG, PG-13
	// - Adulto: todo lo disponible
	switch userAgeRating {
	case "Niño":
		query = `
			SELECT id, title, type, genre, duration, age_rating, synopsis, release_year, director, average_rating, is_available
			FROM audiovisual_content
			WHERE age_rating IN ('G') AND is_available = 1
			ORDER BY average_rating DESC
		`
	case "Adolescente":
		query = `
			SELECT id, title, type, genre, duration, age_rating, synopsis, release_year, director, average_rating, is_available
			FROM audiovisual_content
			WHERE age_rating IN ('G', 'PG', 'PG-13') AND is_available = 1
			ORDER BY average_rating DESC
		`
	default: // Adulto u otros
		query = `
			SELECT id, title, type, genre, duration, age_rating, synopsis, release_year, director, average_rating, is_available
			FROM audiovisual_content
			WHERE is_available = 1
			ORDER BY average_rating DESC
		`
	}

	rows, err := conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contents []models.AudiovisualContent
	for rows.Next() {
		var c models.AudiovisualContent
		err := rows.Scan(
			&c.ID,
			&c.Title,
			&c.Type,
			&c.Genre,
			&c.Duration,
			&c.AgeRating,
			&c.Synopsis,
			&c.ReleaseYear,
			&c.Director,
			&c.AverageRating,
			&c.IsAvailable,
		)
		if err != nil {
			return nil, err
		}
		contents = append(contents, c)
	}

	return contents, nil
}

// --- AUDIO ---

func (r *sqliteContentRepo) CreateAudio(content *models.AudioContent) error {
	conn := db.GetDB()

	query := `
		INSERT INTO audio_content (title, type, genre, duration, age_rating, artist, album, track_number, average_rating, is_available)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := conn.Exec(query,
		content.Title,
		content.Type,
		content.Genre,
		content.Duration,
		content.AgeRating,
		content.Artist,
		content.Album,
		content.TrackNumber,
		content.AverageRating,
		content.IsAvailable,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	content.ID = int(id)
	return nil
}

func (r *sqliteContentRepo) FindAudioByID(id int) (*models.AudioContent, error) {
	conn := db.GetDB()

	query := `
		SELECT id, title, type, genre, duration, age_rating, artist, album, track_number, average_rating, is_available
		FROM audio_content
		WHERE id = ?
	`

	var c models.AudioContent
	err := conn.QueryRow(query, id).Scan(
		&c.ID,
		&c.Title,
		&c.Type,
		&c.Genre,
		&c.Duration,
		&c.AgeRating,
		&c.Artist,
		&c.Album,
		&c.TrackNumber,
		&c.AverageRating,
		&c.IsAvailable,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("contenido de audio no encontrado")
	}
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *sqliteContentRepo) FindAllAudio() ([]models.AudioContent, error) {
	conn := db.GetDB()

	query := `
		SELECT id, title, type, genre, duration, age_rating, artist, album, track_number, average_rating, is_available
		FROM audio_content
		WHERE is_available = 1
		ORDER BY average_rating DESC
	`

	rows, err := conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contents []models.AudioContent
	for rows.Next() {
		var c models.AudioContent
		err := rows.Scan(
			&c.ID,
			&c.Title,
			&c.Type,
			&c.Genre,
			&c.Duration,
			&c.AgeRating,
			&c.Artist,
			&c.Album,
			&c.TrackNumber,
			&c.AverageRating,
			&c.IsAvailable,
		)
		if err != nil {
			return nil, err
		}
		contents = append(contents, c)
	}

	return contents, nil
}

func (r *sqliteContentRepo) SearchAudioByTitle(title string) ([]models.AudioContent, error) {
	conn := db.GetDB()

	query := `
		SELECT id, title, type, genre, duration, age_rating, artist, album, track_number, average_rating, is_available
		FROM audio_content
		WHERE title LIKE ? AND is_available = 1
		ORDER BY average_rating DESC
	`

	rows, err := conn.Query(query, "%"+title+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contents []models.AudioContent
	for rows.Next() {
		var c models.AudioContent
		err := rows.Scan(
			&c.ID,
			&c.Title,
			&c.Type,
			&c.Genre,
			&c.Duration,
			&c.AgeRating,
			&c.Artist,
			&c.Album,
			&c.TrackNumber,
			&c.AverageRating,
			&c.IsAvailable,
		)
		if err != nil {
			return nil, err
		}
		contents = append(contents, c)
	}

	return contents, nil
}

func (r *sqliteContentRepo) FindAllAudioAllowed(userAgeRating string) ([]models.AudioContent, error) {
	conn := db.GetDB()

	var query string

	// Regla simple:
	// - Niño, Adolescente: solo General
	// - Adulto: General y Explicit
	switch userAgeRating {
	case "Niño", "Adolescente":
		query = `
			SELECT id, title, type, genre, duration, age_rating, artist, album, track_number, average_rating, is_available
			FROM audio_content
			WHERE age_rating = 'General' AND is_available = 1
			ORDER BY average_rating DESC
		`
	default: // Adulto
		query = `
			SELECT id, title, type, genre, duration, age_rating, artist, album, track_number, average_rating, is_available
			FROM audio_content
			WHERE is_available = 1
			ORDER BY average_rating DESC
		`
	}

	rows, err := conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contents []models.AudioContent
	for rows.Next() {
		var c models.AudioContent
		err := rows.Scan(
			&c.ID,
			&c.Title,
			&c.Type,
			&c.Genre,
			&c.Duration,
			&c.AgeRating,
			&c.Artist,
			&c.Album,
			&c.TrackNumber,
			&c.AverageRating,
			&c.IsAvailable,
		)
		if err != nil {
			return nil, err
		}
		contents = append(contents, c)
	}

	return contents, nil
}

// --- RATINGS ---

func (r *sqliteContentRepo) UpdateAverageRating(contentID int, contentType string, avg float64) error {
	conn := db.GetDB()

	var query string
	if contentType == "audiovisual" {
		query = `UPDATE audiovisual_content SET average_rating = ? WHERE id = ?`
	} else {
		query = `UPDATE audio_content SET average_rating = ? WHERE id = ?`
	}

	_, err := conn.Exec(query, avg, contentID)
	return err
}
