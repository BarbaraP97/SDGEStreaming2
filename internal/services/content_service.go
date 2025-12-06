// internal/services/content_service.go
package services

import (
	"SDGEStreaming/internal/models"
	"SDGEStreaming/internal/repositories"
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

func (s *ContentService) SearchAudioByTitle(title string) ([]models.AudioContent, error) {
	return s.contentRepo.SearchAudioByTitle(title)
}
