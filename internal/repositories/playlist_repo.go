// internal/repositories/playlist_repo.go
package repositories

import (
	"SDGEStreaming/internal/db"
	"SDGEStreaming/internal/models"
	"database/sql"
)

type PlaylistRepo interface {
	CreatePlaylist(p *models.Playlist) error
	FindPlaylistsByProfileID(profileID int) ([]models.Playlist, error)
	DeletePlaylist(id int) error
	AddItem(item *models.PlaylistItem) error
	RemoveItem(playlistID, contentID int, contentType string) error
	GetItems(playlistID int) ([]models.PlaylistItem, error)
}

type sqlitePlaylistRepo struct{}

func NewPlaylistRepo() PlaylistRepo {
	return &sqlitePlaylistRepo{}
}

func (r *sqlitePlaylistRepo) CreatePlaylist(p *models.Playlist) error {
	conn := db.GetDB()
	query := `INSERT INTO playlists (profile_id, name, description, created_at) VALUES (?, ?, ?, ?)`
	result, err := conn.Exec(query, p.ProfileID, p.Name, p.Description, p.CreatedAt)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	p.ID = int(id)
	return nil
}

func (r *sqlitePlaylistRepo) FindPlaylistsByProfileID(profileID int) ([]models.Playlist, error) {
	conn := db.GetDB()
	query := `SELECT id, profile_id, name, description, created_at FROM playlists WHERE profile_id = ? ORDER BY created_at DESC`
	rows, err := conn.Query(query, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playlists []models.Playlist
	for rows.Next() {
		var p models.Playlist
		if err := rows.Scan(&p.ID, &p.ProfileID, &p.Name, &p.Description, &p.CreatedAt); err != nil {
			return nil, err
		}
		playlists = append(playlists, p)
	}
	return playlists, nil
}

func (r *sqlitePlaylistRepo) DeletePlaylist(id int) error {
	conn := db.GetDB()
	_, err := conn.Exec("DELETE FROM playlists WHERE id = ?", id)
	return err
}

func (r *sqlitePlaylistRepo) AddItem(item *models.PlaylistItem) error {
	conn := db.GetDB()
	query := `INSERT INTO playlist_items (playlist_id, content_id, content_type, position) VALUES (?, ?, ?, ?)`
	result, err := conn.Exec(query, item.PlaylistID, item.ContentID, item.ContentType, item.Position)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	item.ID = int(id)
	return nil
}

func (r *sqlitePlaylistRepo) RemoveItem(playlistID, contentID int, contentType string) error {
	conn := db.GetDB()
	_, err := conn.Exec("DELETE FROM playlist_items WHERE playlist_id = ? AND content_id = ? AND content_type = ?", 
		playlistID, contentID, contentType)
	return err
}

func (r *sqlitePlaylistRepo) GetItems(playlistID int) ([]models.PlaylistItem, error) {
	conn := db.GetDB()
	query := `SELECT id, playlist_id, content_id, content_type, position FROM playlist_items WHERE playlist_id = ? ORDER BY position`
	rows, err := conn.Query(query, playlistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.PlaylistItem
	for rows.Next() {
		var item models.PlaylistItem
		if err := rows.Scan(&item.ID, &item.PlaylistID, &item.ContentID, &item.ContentType, &item.Position); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}