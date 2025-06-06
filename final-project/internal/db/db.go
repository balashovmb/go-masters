package db

import (
	"context"
	"go-masters/final-project/internal/models"
)

type DB interface {
	AddAlbum(context.Context, models.Album) error
	ListAlbums(context.Context) ([]models.Album, error)
	Test(context.Context) ([]models.User, error)
	AddReview(context.Context, models.Review) error
}
