// internal/services/content_service.go
package services

import (
	"SDGEStreaming/internal/db"
	"SDGEStreaming/internal/models"
	"SDGEStreaming/internal/repositories"
	"fmt"
)

// ContentService handles content-related business logic.
type ContentService struct {
	contentRepo repositories.ContentRepo
}

func NewContentService(contentRepo repositories.ContentRepo) *ContentService {
	return &ContentService{contentRepo: contentRepo}
}

// --- AUDIOVISUAL ---
func (s *ContentService) CreateAudiovisual(title, contentType, genre string, duration int, ageRating, synopsis string, releaseYear int, director string) error {
	content := &models.AudiovisualContent{
		Title:       title,
		Type:        contentType,
		Genre:       genre,
		Duration:    duration,
		AgeRating:   ageRating,
		Synopsis:    synopsis,
		ReleaseYear: releaseYear,
		Director:    director,
	}
	return s.contentRepo.CreateAudiovisual(content)
}

func (s *ContentService) GetAudiovisualByID(id int) (*models.AudiovisualContent, error) {
	return s.contentRepo.FindAudiovisualByID(id)
}

func (s *ContentService) GetAllAudiovisual() ([]models.AudiovisualContent, error) {
	return s.contentRepo.FindAllAudiovisual()
}

func (s *ContentService) GetAllAudiovisualForUser(ageRating string) ([]models.AudiovisualContent, error) {
	return s.contentRepo.FindAllAudiovisualAllowed(ageRating)
}

func (s *ContentService) SearchAudiovisualByTitle(title string) ([]models.AudiovisualContent, error) {
	return s.contentRepo.SearchAudiovisualByTitle(title)
}

// --- AUDIO ---
func (s *ContentService) CreateAudio(title, contentType, genre string, duration int, ageRating, artist, album string, trackNumber int) error {
	content := &models.AudioContent{
		Title:       title,
		Type:        contentType,
		Genre:       genre,
		Duration:    duration,
		AgeRating:   ageRating,
		Artist:      artist,
		Album:       album,
		TrackNumber: trackNumber,
	}
	return s.contentRepo.CreateAudio(content)
}

func (s *ContentService) GetAudioByID(id int) (*models.AudioContent, error) {
	return s.contentRepo.FindAudioByID(id)
}

func (s *ContentService) GetAllAudio() ([]models.AudioContent, error) {
	return s.contentRepo.FindAllAudio()
}

func (s *ContentService) GetAllAudioForUser(ageRating string) ([]models.AudioContent, error) {
	return s.contentRepo.FindAllAudioAllowed(ageRating)
}

func (s *ContentService) SearchAudioByTitle(title string) ([]models.AudioContent, error) {
	return s.contentRepo.SearchAudioByTitle(title)
}

// --- CALIFICACIONES ---
func (s *ContentService) RateContent(userID, contentID int, contentType string, rating float64) error {
	if rating < 1.0 || rating > 10.0 {
		return fmt.Errorf("la calificación debe estar entre 1.0 y 10.0")
	}

	conn := db.GetDB()

	// Insertar o actualizar calificación
	_, err := conn.Exec(`
		INSERT INTO user_ratings (user_id, content_id, content_type, rating)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(user_id, content_id, content_type) 
		DO UPDATE SET rating = ?, rated_at = CURRENT_TIMESTAMP
	`, userID, contentID, contentType, rating, rating)

	if err != nil {
		return fmt.Errorf("error al guardar calificación: %w", err)
	}

	// Recalcular promedio
	return s.updateAverageRating(contentID, contentType)
}

func (s *ContentService) updateAverageRating(contentID int, contentType string) error {
	conn := db.GetDB()
	var avgRating float64

	err := conn.QueryRow(`
		SELECT COALESCE(AVG(rating), 0.0)
		FROM user_ratings
		WHERE content_id = ? AND content_type = ?
	`, contentID, contentType).Scan(&avgRating)

	if err != nil {
		return fmt.Errorf("error al calcular promedio: %w", err)
	}

	return s.contentRepo.UpdateAverageRating(contentID, contentType, avgRating)
}
