package db

import (
	"context"
	"go-masters/final-project/internal/models"
)

type DB interface {
	AddAlbum(context.Context, models.Album) error
	AddReview(context.Context, models.Review) (int, error)
	ListReviews(context.Context, string, int) ([]models.Review, error)
	UpdateReviewRating(context.Context, int, int) error
	AverageRating(context.Context, int) (float64, error)
}
