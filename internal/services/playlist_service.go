// internal/services/playlist_service.go
package services

import (
	"SDGEStreaming/internal/models"
	"SDGEStreaming/internal/repositories"
	"time"
)

type PlaylistService struct {
	playlistRepo repositories.PlaylistRepo
	contentRepo  repositories.ContentRepo
}

func NewPlaylistService(playlistRepo repositories.PlaylistRepo, contentRepo repositories.ContentRepo) *PlaylistService {
	return &PlaylistService{
		playlistRepo: playlistRepo,
		contentRepo:  contentRepo,
	}
}

func (s *PlaylistService) CreatePlaylist(profileID int, name, description string) (*models.Playlist, error) {
	playlist := &models.Playlist{
		ProfileID:   profileID,
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
	}
	err := s.playlistRepo.CreatePlaylist(playlist)
	return playlist, err
}

func (s *PlaylistService) GetPlaylistsByProfileID(profileID int) ([]models.Playlist, error) {
	return s.playlistRepo.FindPlaylistsByProfileID(profileID)
}

func (s *PlaylistService) DeletePlaylist(id int) error {
	return s.playlistRepo.DeletePlaylist(id)
}

func (s *PlaylistService) AddItemToPlaylist(playlistID, contentID int, contentType string) error {
	items, _ := s.playlistRepo.GetItems(playlistID)
	position := len(items) + 1

	item := &models.PlaylistItem{
		PlaylistID:  playlistID,
		ContentID:   contentID,
		ContentType: contentType,
		Position:    position,
	}
	return s.playlistRepo.AddItem(item)
}

func (s *PlaylistService) RemoveItemFromPlaylist(playlistID, contentID int, contentType string) error {
	return s.playlistRepo.RemoveItem(playlistID, contentID, contentType)
}

func (s *PlaylistService) GetPlaylistItems(playlistID int) ([]models.PlaylistItem, error) {
	return s.playlistRepo.GetItems(playlistID)
}