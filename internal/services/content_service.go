// internal/services/content_service.go
package services

import (
	"SDGEStreaming/internal/db"
	"SDGEStreaming/internal/models"
	"SDGEStreaming/internal/repositories"
	"fmt"
	"strings"
)

type ContentService struct {
	contentRepo repositories.ContentRepo
}

func NewContentService(contentRepo repositories.ContentRepo) *ContentService {
	return &ContentService{contentRepo: contentRepo}
}

// AUDIOVISUAL
func (s *ContentService) CreateAudiovisual(title, contentType, genre string, duration int, ageRating, synopsis string, releaseYear int, director, actors string) error {
	content := &models.AudiovisualContent{
		Title:       title,
		Type:        contentType,
		Genre:       genre,
		Duration:    duration,
		AgeRating:   ageRating,
		Synopsis:    synopsis,
		ReleaseYear: releaseYear,
		Director:    director,
		Actors:      actors,
	}
	return s.contentRepo.CreateAudiovisual(content)
}

func (s *ContentService) GetAudiovisualByID(id int) (*models.AudiovisualContent, error) {
	return s.contentRepo.FindAudiovisualByID(id)
}

func (s *ContentService) GetAllAudiovisual() ([]models.AudiovisualContent, error) {
	return s.contentRepo.FindAllAudiovisual()
}

func (s *ContentService) SearchAudiovisualByTitle(title string) ([]models.AudiovisualContent, error) {
	return s.contentRepo.SearchAudiovisualByTitle(title)
}

// AUDIO
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

func (s *ContentService) SearchAudioByTitle(title string) ([]models.AudioContent, error) {
	return s.contentRepo.SearchAudioByTitle(title)
}

// BUSQUEDA AVANZADA
func (s *ContentService) SearchByTitle(query, ageRating string) []interface{} {
	var results []interface{}
	
	avs, _ := s.contentRepo.SearchAudiovisualByTitle(query)
	for _, av := range avs {
		if isContentAllowed(av.AgeRating, ageRating) {
			results = append(results, av)
		}
	}
	
	audios, _ := s.contentRepo.SearchAudioByTitle(query)
	for _, a := range audios {
		if isContentAllowed(a.AgeRating, ageRating) {
			results = append(results, a)
		}
	}
	
	return results
}

func (s *ContentService) SearchByGenre(genre, ageRating string) []interface{} {
	var results []interface{}
	
	avs, _ := s.contentRepo.FindAllAudiovisual()
	for _, av := range avs {
		if strings.Contains(strings.ToLower(av.Genre), strings.ToLower(genre)) && isContentAllowed(av.AgeRating, ageRating) {
			results = append(results, av)
		}
	}
	
	audios, _ := s.contentRepo.FindAllAudio()
	for _, a := range audios {
		if strings.Contains(strings.ToLower(a.Genre), strings.ToLower(genre)) && isContentAllowed(a.AgeRating, ageRating) {
			results = append(results, a)
		}
	}
	
	return results
}

func (s *ContentService) SearchByActor(actor, ageRating string) []models.AudiovisualContent {
	var results []models.AudiovisualContent
	
	avs, _ := s.contentRepo.FindAllAudiovisual()
	for _, av := range avs {
		if strings.Contains(strings.ToLower(av.Actors), strings.ToLower(actor)) && isContentAllowed(av.AgeRating, ageRating) {
			results = append(results, av)
		}
	}
	
	return results
}

// CALIFICACIONES
func (s *ContentService) RateContent(profileID, contentID int, contentType string, rating float64) error {
	if rating < 1.0 || rating > 10.0 {
		return fmt.Errorf("la calificacion debe estar entre 1.0 y 10.0")
	}
	
	conn := db.GetDB()
	
	_, err := conn.Exec(`
		INSERT INTO user_ratings (profile_id, content_id, content_type, rating)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(profile_id, content_id, content_type) 
		DO UPDATE SET rating = ?, rated_at = CURRENT_TIMESTAMP
	`, profileID, contentID, contentType, rating, rating)
	
	if err != nil {
		return fmt.Errorf("error al guardar calificacion: %w", err)
	}
	
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

func isContentAllowed(contentRating, profileRating string) bool {
	ratings := map[string]int{"G": 1, "PG": 2, "PG-13": 3, "R": 4, "General": 1, "Explicit": 4}
	contentLevel := ratings[contentRating]
	profileLevel := ratings[profileRating]
	return contentLevel <= profileLevel
}