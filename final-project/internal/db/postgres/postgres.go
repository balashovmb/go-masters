package postgres

import (
	"context"
	"fmt"
	"go-masters/final-project/internal/models"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func New(connstr string) (*Postgres, error) {
	pool, err := pgxpool.New(context.Background(), connstr)
	if err != nil {
		return nil, err
	}

	pg := Postgres{pool: pool}

	err = pg.applyMigrations()
	if err != nil {
		return nil, err
	}

	return &pg, nil
}

func (pg *Postgres) applyMigrations() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	parent := filepath.Dir(wd)
	migrationsFS := os.DirFS(parent)
	goose.SetBaseFS(migrationsFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	db := stdlib.OpenDBFromPool(pg.pool)
	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}
	return db.Close()
}

func (pg *Postgres) AddAlbum(ctx context.Context, album models.Album) error {
	_, err := pg.pool.Exec(
		ctx,
		"INSERT INTO albums (id  artist, title, year) VALUES ($1, $2, $3, $4)",
		album.ID,
		album.Artist,
		album.Title,
		album.Year)
	return err
}

func (pg *Postgres) AddReview(ctx context.Context, review models.Review) (int, error) {
	var id int
	err := pg.pool.QueryRow(
		ctx,
		`INSERT INTO reviews (object_id, user_id, text, rating)
         VALUES ($1, $2, $3, $4)
         RETURNING id`,
		review.ObjectID,
		review.UserID,
		review.Text,
		review.Rating,
	).Scan(&id)
	return id, err
}

func (pg *Postgres) UpdateReviewRating(ctx context.Context, id int, rating int) error {
	_, err := pg.pool.Exec(
		ctx,
		"UPDATE reviews SET rating = $1 WHERE id = $2",
		rating,
		id,
	)
	return err
}

func (pg *Postgres) ListReviews(ctx context.Context, filter string, id int) ([]models.Review, error) {
	log.Debug().Str("filter", filter).Int("id", id).Msg("ListReviews")
	rows, err := pg.pool.Query(ctx, fmt.Sprintf("SELECT id, object_id, user_id, text, rating FROM reviews WHERE %s = $1", filter), id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var review models.Review

		if err := rows.Scan(
			&review.ID,
			&review.ObjectID,
			&review.UserID,
			&review.Text,
			&review.Rating,
		); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		reviews = append(reviews, review)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}

func (p *Postgres) AverageRating(ctx context.Context, id int) (float64, error) {
	var rating float64
	err := p.pool.QueryRow(ctx, "SELECT AVG(rating) FROM reviews WHERE object_id = $1 AND rating != 0", id).Scan(&rating)

	return rating, err
}
