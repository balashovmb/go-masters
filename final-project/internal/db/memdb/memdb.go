package memdb

import (
	"context"

	"go-masters/final-project/internal/models"
)

type MemDB struct {
	data []models.Album
}

func New() *MemDB {
	return &MemDB{}
}

func (m *MemDB) AddAlbum(_ context.Context, album models.Album) error {
	m.data = append(m.data, album)
	return nil
}

func (m *MemDB) ListAlbums(_ context.Context) ([]models.Album, error) {
	return m.data, nil
}
