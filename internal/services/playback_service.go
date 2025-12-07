// internal/services/playback_service.go
// Gestiona el historial de reproducción, favoritos y progreso del usuario.
package services

import (
	"SDGEStreaming/internal/models"
	"SDGEStreaming/internal/repositories"
	"fmt"
)

// PlaybackService encapsula la lógica de negocio para la reproducción.
type PlaybackService struct {
	historyRepo  repositories.PlaybackHistoryRepo
	favoriteRepo repositories.FavoriteRepo
	contentRepo  repositories.ContentRepo
}

// NewPlaybackService crea una nueva instancia del servicio.
func NewPlaybackService(historyRepo repositories.PlaybackHistoryRepo, favoriteRepo repositories.FavoriteRepo, contentRepo repositories.ContentRepo) *PlaybackService {
	return &PlaybackService{
		historyRepo:  historyRepo,
		favoriteRepo: favoriteRepo,
		contentRepo:  contentRepo,
	}
}

// AddToHistory agrega una entrada al historial de reproducción.
func (s *PlaybackService) AddToHistory(userID, contentID int, contentType string) error {
	if contentType != "audio" && contentType != "audiovisual" {
		return fmt.Errorf("tipo de contenido inválido")
	}

	// Verificar que el contenido exista
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
		UserID:      userID,
		ContentID:   contentID,
		ContentType: contentType,
		Progress:    0, // Se puede actualizar más tarde
	}
	return s.historyRepo.Create(entry)
}

// UpdateProgress actualiza el progreso de reproducción de un contenido.
func (s *PlaybackService) UpdateProgress(userID, contentID int, contentType string, progressSeconds int) error {
	if progressSeconds < 0 {
		return fmt.Errorf("el progreso no puede ser negativo")
	}

	return s.historyRepo.UpdateProgress(userID, contentID, contentType, progressSeconds)
}

// GetHistory obtiene el historial de reproducción de un usuario (últimas 10 entradas).
func (s *PlaybackService) GetHistory(userID int) ([]models.PlaybackHistory, error) {
	return s.historyRepo.FindByUserID(userID)
}

// AddFavorite agrega un contenido a la lista de favoritos del usuario.
func (s *PlaybackService) AddFavorite(userID, contentID int, contentType string) error {
	if contentType != "audio" && contentType != "audiovisual" {
		return fmt.Errorf("tipo de contenido inválido para favoritos")
	}

	// Verificar que el contenido exista (mismo código que en AddToHistory)
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
		UserID:      userID,
		ContentID:   contentID,
		ContentType: contentType,
	}
	return s.favoriteRepo.Create(favorite)
}

// RemoveFavorite elimina un contenido de la lista de favoritos.
func (s *PlaybackService) RemoveFavorite(userID, contentID int, contentType string) error {
	return s.favoriteRepo.Delete(userID, contentID, contentType)
}

// GetFavorites obtiene la lista de favoritos de un usuario.
func (s *PlaybackService) GetFavorites(userID int) ([]models.Favorite, error) {
	return s.favoriteRepo.FindByUserID(userID)
}

// GetContinueWatching obtiene los contenidos donde el usuario dejó de ver/escuchar.
// Devuelve los últimos 5 elementos con progreso > 0.
func (s *PlaybackService) GetContinueWatching(userID int) ([]models.PlaybackHistory, error) {
	return s.historyRepo.FindContinueWatching(userID)
}

// GetRecommendations genera recomendaciones simples basadas en el género de los favoritos.
// Este es un ejemplo básico; en un sistema real se usaría un algoritmo más complejo.
func (s *PlaybackService) GetRecommendations(userID int) ([]interface{}, error) {
	favorites, err := s.GetFavorites(userID)
	if err != nil {
		return nil, err
	}

	if len(favorites) == 0 {
		// Si no hay favoritos, devolver contenido popular (los primeros 5)
		audiovisuals, err := s.contentRepo.FindAllAudiovisual()
if err != nil {
    // log o devolver error
    return nil, fmt.Errorf("no se pudo obtener contenido audiovisual: %w", err)
}
audios, err := s.contentRepo.FindAllAudio()
if err != nil {
    return nil, fmt.Errorf("no se pudo obtener contenido de audio: %w", err)
}
		var recommendations []interface{}
		for i, av := range audiovisuals {
			if i >= 3 {
				break
			}
			recommendations = append(recommendations, av)
		}
		for i, a := range audios {
			if i >= 2 {
				break
			}
			recommendations = append(recommendations, a)
		}
		return recommendations, nil
	}

	// Contar géneros de los favoritos
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

	// Encontrar el género más popular
	mostPopularGenre := ""
	maxCount := 0
	for genre, count := range genreCount {
		if count > maxCount {
			maxCount = count
			mostPopularGenre = genre
		}
	}

	// Buscar contenido en ese género
	var recommendations []interface{}
	if mostPopularGenre != "" {
		audiovisuals, _ := s.contentRepo.SearchAudiovisualByTitle(mostPopularGenre)
		audios, _ := s.contentRepo.SearchAudioByTitle(mostPopularGenre)
		for i, av := range audiovisuals {
			if i >= 3 {
				break
			}
			recommendations = append(recommendations, av)
		}
		for i, a := range audios {
			if i >= 2 {
				break
			}
			recommendations = append(recommendations, a)
		}
	}

	return recommendations, nil
}
