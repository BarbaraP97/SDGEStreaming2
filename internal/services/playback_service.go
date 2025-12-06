// internal/services/playback_service.go
package services

import (
	"SDGEStreaming/internal/models"
	"SDGEStreaming/internal/repositories"
	"fmt"
	"strings"
)

type PlaybackService struct {
	historyRepo  repositories.PlaybackHistoryRepo
	favoriteRepo repositories.FavoriteRepo
	contentRepo  repositories.ContentRepo
}

func NewPlaybackService(historyRepo repositories.PlaybackHistoryRepo, favoriteRepo repositories.FavoriteRepo, contentRepo repositories.ContentRepo) *PlaybackService {
	return &PlaybackService{
		historyRepo:  historyRepo,
		favoriteRepo: favoriteRepo,
		contentRepo:  contentRepo,
	}
}

func (s *PlaybackService) AddToHistory(profileID, contentID int, contentType string) error {
	if contentType != "audio" && contentType != "audiovisual" {
		return fmt.Errorf("tipo de contenido invalido")
	}

	if contentType == "audiovisual" {
		_, err := s.contentRepo.FindAudiovisualByID(contentID)
		if err != nil {
			return fmt.Errorf("contenido audiovisual no encontrado")
		}
	} else {
		_, err := s.contentRepo.FindAudioByID(contentID)
		if err != nil {
			return fmt.Errorf("contenido de audio no encontrado")
		}
	}

	entry := &models.PlaybackHistory{
		ProfileID:   profileID,
		ContentID:   contentID,
		ContentType: contentType,
		Progress:    0,
	}
	return s.historyRepo.Create(entry)
}

func (s *PlaybackService) UpdateProgress(profileID, contentID int, contentType string, progressSeconds int) error {
	if progressSeconds < 0 {
		return fmt.Errorf("el progreso no puede ser negativo")
	}

	return s.historyRepo.UpdateProgress(profileID, contentID, contentType, progressSeconds)
}

func (s *PlaybackService) GetHistory(profileID int) ([]models.PlaybackHistory, error) {
	return s.historyRepo.FindByProfileID(profileID)
}

func (s *PlaybackService) AddFavorite(profileID, contentID int, contentType string) error {
	if contentType != "audio" && contentType != "audiovisual" {
		return fmt.Errorf("tipo de contenido invalido para favoritos")
	}

	if contentType == "audiovisual" {
		_, err := s.contentRepo.FindAudiovisualByID(contentID)
		if err != nil {
			return fmt.Errorf("contenido audiovisual no encontrado")
		}
	} else {
		_, err := s.contentRepo.FindAudioByID(contentID)
		if err != nil {
			return fmt.Errorf("contenido de audio no encontrado")
		}
	}

	favorite := &models.Favorite{
		ProfileID:   profileID,
		ContentID:   contentID,
		ContentType: contentType,
	}
	return s.favoriteRepo.Create(favorite)
}

func (s *PlaybackService) RemoveFavorite(profileID, contentID int, contentType string) error {
	return s.favoriteRepo.Delete(profileID, contentID, contentType)
}

func (s *PlaybackService) GetFavorites(profileID int) ([]models.Favorite, error) {
	return s.favoriteRepo.FindByProfileID(profileID)
}

func (s *PlaybackService) GetContinueWatching(profileID int) ([]models.PlaybackHistory, error) {
	return s.historyRepo.FindContinueWatching(profileID)
}

func (s *PlaybackService) GetRecommendations(profileID int, ageRating string) []interface{} {
	var recommendations []interface{}
	
	favorites, _ := s.GetFavorites(profileID)
	history, _ := s.GetHistory(profileID)
	
	if len(favorites) == 0 && len(history) == 0 {
		audiovisuals, _ := s.contentRepo.FindAllAudiovisual()
		for i, av := range audiovisuals {
			if i >= 5 {
				break
			}
			if isAllowed(av.AgeRating, ageRating) {
				recommendations = append(recommendations, av)
			}
		}
		return recommendations
	}

	genreCount := make(map[string]int)
	for _, fav := range favorites {
		var genre string
		if fav.ContentType == "audiovisual" {
			content, err := s.contentRepo.FindAudiovisualByID(fav.ContentID)
			if err != nil {
				continue
			}
			genre = content.Genre
		} else {
			content, err := s.contentRepo.FindAudioByID(fav.ContentID)
			if err != nil {
				continue
			}
			genre = content.Genre
		}
		genreCount[genre]++
	}

	mostPopularGenre := ""
	maxCount := 0
	for genre, count := range genreCount {
		if count > maxCount {
			maxCount = count
			mostPopularGenre = genre
		}
	}

	if mostPopularGenre != "" {
		audiovisuals, _ := s.contentRepo.FindAllAudiovisual()
		for _, av := range audiovisuals {
			if strings.Contains(strings.ToLower(av.Genre), strings.ToLower(mostPopularGenre)) && isAllowed(av.AgeRating, ageRating) {
				recommendations = append(recommendations, av)
				if len(recommendations) >= 10 {
					break
				}
			}
		}
	}

	return recommendations
}

func isAllowed(contentRating, profileRating string) bool {
	ratings := map[string]int{"G": 1, "PG": 2, "PG-13": 3, "R": 4, "General": 1, "Explicit": 4}
	contentLevel := ratings[contentRating]
	profileLevel := ratings[profileRating]
	return contentLevel <= profileLevel
}