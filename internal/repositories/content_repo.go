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

	CreateAudio(content *models.AudioContent) error
	FindAudioByID(id int) (*models.AudioContent, error)
	FindAllAudio() ([]models.AudioContent, error)
	SearchAudioByTitle(title string) ([]models.AudioContent, error)

	// Compartido
	UpdateAvailability(contentID int, contentType string, available bool) error
	UpdateAverageRating(contentID int, contentType string, rating float64) error
	GetDB() *sql.DB
}

type sqliteContentRepo struct {
	conn *sql.DB
}

func NewContentRepo() ContentRepo {
	return &sqliteContentRepo{
		conn: db.GetDB(),
	}
}

func (r *sqliteContentRepo) GetDB() *sql.DB {
	return r.conn
}

//
// AUDIOVISUAL
//

func (r *sqliteContentRepo) CreateAudiovisual(c *models.AudiovisualContent) error {
	query := `
		INSERT INTO audiovisual_content
		(title, type, genre, duration, age_rating, synopsis, release_year, director)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.conn.Exec(query,
		c.Title, c.Type, c.Genre, c.Duration, c.AgeRating,
		c.Synopsis, c.ReleaseYear, c.Director,
	)
	if err != nil {
		return fmt.Errorf("error creating audiovisual content: %w", err)
	}

	return nil
}

func (r *sqliteContentRepo) FindAudiovisualByID(id int) (*models.AudiovisualContent, error) {
	query := `
		SELECT id, title, type, genre, duration, age_rating,
		       synopsis, release_year, director,
		       average_rating, is_available
		FROM audiovisual_content
		WHERE id = ?
	`

	row := r.conn.QueryRow(query, id)

	var c models.AudiovisualContent
	err := row.Scan(
		&c.ID, &c.Title, &c.Type, &c.Genre, &c.Duration, &c.AgeRating,
		&c.Synopsis, &c.ReleaseYear, &c.Director,
		&c.AverageRating, &c.IsAvailable,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("error scanning audiovisual content: %w", err)
	}

	return &c, nil
}

func (r *sqliteContentRepo) FindAllAudiovisual() ([]models.AudiovisualContent, error) {
	query := `
		SELECT id, title, type, genre, duration, age_rating,
			   synopsis, release_year, director,
			   average_rating, is_available
		FROM audiovisual_content
		ORDER BY release_year DESC
	`

	rows, err := r.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error fetching audiovisual content: %w", err)
	}
	defer rows.Close()

	var list []models.AudiovisualContent

	for rows.Next() {
		var c models.AudiovisualContent
		if err := rows.Scan(
			&c.ID, &c.Title, &c.Type, &c.Genre, &c.Duration, &c.AgeRating,
			&c.Synopsis, &c.ReleaseYear, &c.Director,
			&c.AverageRating, &c.IsAvailable,
		); err != nil {
			return nil, fmt.Errorf("error scanning audiovisual row: %w", err)
		}
		list = append(list, c)
	}

	return list, nil
}

func (r *sqliteContentRepo) SearchAudiovisualByTitle(title string) ([]models.AudiovisualContent, error) {
	query := `
		SELECT id, title, type, genre, duration, age_rating,
			   synopsis, release_year, director,
			   average_rating, is_available
		FROM audiovisual_content
		WHERE title LIKE ?
		ORDER BY release_year DESC
	`
	rows, err := r.conn.Query(query, "%"+title+"%")
	if err != nil {
		return nil, fmt.Errorf("error searching audiovisual content by title: %w", err)
	}
	defer rows.Close()

	var list []models.AudiovisualContent
	for rows.Next() {
		var c models.AudiovisualContent
		if err := rows.Scan(
			&c.ID, &c.Title, &c.Type, &c.Genre, &c.Duration, &c.AgeRating,
			&c.Synopsis, &c.ReleaseYear, &c.Director,
			&c.AverageRating, &c.IsAvailable,
		); err != nil {
			return nil, fmt.Errorf("error scanning audiovisual row: %w", err)
		}
		list = append(list, c)
	}
	return list, nil
}

//
// AUDIO
//

func (r *sqliteContentRepo) CreateAudio(c *models.AudioContent) error {
	query := `
		INSERT INTO audio_content
		(title, type, genre, duration, age_rating, artist, album, track_number)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.conn.Exec(query,
		c.Title, c.Type, c.Genre, c.Duration, c.AgeRating,
		c.Artist, c.Album, c.TrackNumber,
	)
	if err != nil {
		return fmt.Errorf("error creating audio content: %w", err)
	}

	return nil
}

func (r *sqliteContentRepo) FindAudioByID(id int) (*models.AudioContent, error) {
	query := `
		SELECT id, title, type, genre, duration, age_rating,
		       artist, album, track_number,
		       average_rating, is_available
		FROM audio_content
		WHERE id = ?
	`

	row := r.conn.QueryRow(query, id)

	var c models.AudioContent
	err := row.Scan(
		&c.ID, &c.Title, &c.Type, &c.Genre, &c.Duration, &c.AgeRating,
		&c.Artist, &c.Album, &c.TrackNumber,
		&c.AverageRating, &c.IsAvailable,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("error scanning audio content: %w", err)
	}

	return &c, nil
}

func (r *sqliteContentRepo) FindAllAudio() ([]models.AudioContent, error) {
	query := `
		SELECT id, title, type, genre, duration, age_rating,
		       artist, album, track_number,
		       average_rating, is_available
		FROM audio_content
		ORDER BY id DESC
	`

	rows, err := r.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error fetching audio content: %w", err)
	}
	defer rows.Close()

	var list []models.AudioContent

	for rows.Next() {
		var c models.AudioContent
		if err := rows.Scan(
			&c.ID, &c.Title, &c.Type, &c.Genre, &c.Duration, &c.AgeRating,
			&c.Artist, &c.Album, &c.TrackNumber,
			&c.AverageRating, &c.IsAvailable,
		); err != nil {
			return nil, fmt.Errorf("error scanning audio row: %w", err)
		}
		list = append(list, c)
	}

	return list, nil
}

func (r *sqliteContentRepo) SearchAudioByTitle(title string) ([]models.AudioContent, error) {
	query := `
		SELECT id, title, type, genre, duration, age_rating,
			   artist, album, track_number,
			   average_rating, is_available
		FROM audio_content
		WHERE title LIKE ?
		ORDER BY id DESC
	`
	rows, err := r.conn.Query(query, "%"+title+"%")
	if err != nil {
		return nil, fmt.Errorf("error searching audio content by title: %w", err)
	}
	defer rows.Close()

	var list []models.AudioContent
	for rows.Next() {
		var c models.AudioContent
		if err := rows.Scan(
			&c.ID, &c.Title, &c.Type, &c.Genre, &c.Duration, &c.AgeRating,
			&c.Artist, &c.Album, &c.TrackNumber,
			&c.AverageRating, &c.IsAvailable,
		); err != nil {
			return nil, fmt.Errorf("error scanning audio row: %w", err)
		}
		list = append(list, c)
	}
	return list, nil
}

//
// COMPARTIDO (audiovisual | audio)
//

func (r *sqliteContentRepo) UpdateAvailability(id int, contentType string, available bool) error {
	var table string

	if contentType == "audiovisual" {
		table = "audiovisual_content"
	} else if contentType == "audio" {
		table = "audio_content"
	} else {
		return fmt.Errorf("invalid content type")
	}

	query := fmt.Sprintf(`
		UPDATE %s
		SET is_available = ?
		WHERE id = ?
	`, table)

	_, err := r.conn.Exec(query, available, id)
	if err != nil {
		return fmt.Errorf("error updating availability: %w", err)
	}

	return nil
}

func (r *sqliteContentRepo) UpdateAverageRating(id int, contentType string, rating float64) error {
	var table string

	if contentType == "audiovisual" {
		table = "audiovisual_content"
	} else if contentType == "audio" {
		table = "audio_content"
	} else {
		return fmt.Errorf("invalid content type")
	}

	query := fmt.Sprintf(`
		UPDATE %s
		SET average_rating = ?
		WHERE id = ?
	`, table)

	_, err := r.conn.Exec(query, rating, id)
	if err != nil {
		return fmt.Errorf("error updating rating: %w", err)
	}

	return nil
}