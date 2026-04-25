package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"shortener/internal/model"
)

var ErrURLNotFound = errors.New("url not found")

type URLRepository struct {
	db *sql.DB
}

func NewURLRepository(db *sql.DB) *URLRepository {
	return &URLRepository{db: db}
}

func (r *URLRepository) Create(ctx context.Context, url string, shortCode string) (*model.UrlResponse, error) {
	now := time.Now().UTC()
	item := &model.UrlResponse{
		Id:        uuid.New(),
		Url:       url,
		ShortCode: shortCode,
		CreatedAt: now,
		UpdatedAt: now,
	}

	const query = `
INSERT INTO urls (id, url, short_code, created_at, updated_at)
VALUES (?, ?, ?, ?, ?);
`

	if _, err := r.db.ExecContext(ctx, query, item.Id.String(), item.Url, item.ShortCode, item.CreatedAt, item.UpdatedAt); err != nil {
		return nil, fmt.Errorf("create url: %w", err)
	}

	return item, nil
}

func (r *URLRepository) GetByShortCode(ctx context.Context, shortCode string) (*model.UrlResponse, error) {
	const query = `
SELECT id, url, short_code, created_at, updated_at
FROM urls
WHERE short_code = ?;
`

	var (
		id        string
		item      model.UrlResponse
		createdAt time.Time
		updatedAt time.Time
	)

	err := r.db.QueryRowContext(ctx, query, shortCode).Scan(
		&id,
		&item.Url,
		&item.ShortCode,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrURLNotFound
		}
		return nil, fmt.Errorf("get url by short code: %w", err)
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("parse url id: %w", err)
	}

	item.Id = parsedID
	item.CreatedAt = createdAt
	item.UpdatedAt = updatedAt

	return &item, nil
}
